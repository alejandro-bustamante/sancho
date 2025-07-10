package service

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"encoding/json"
	"errors"
	"strings"

	model "github.com/alejandro-bustamante/sancho/server/internal/model"
	db "github.com/alejandro-bustamante/sancho/server/internal/repository"
	"github.com/google/uuid"
)

type deezerTrackSampleResponse struct {
	Preview string `json:"preview"`
}

// Code to handle download's status
type DownloadTracker struct {
	sync.Mutex
	status map[string]model.DownloadStatus
	errs   map[string]string
}

func NewDownloadTracker() *DownloadTracker {
	return &DownloadTracker{
		status: make(map[string]model.DownloadStatus),
		errs:   make(map[string]string),
	}
}

func (dt *DownloadTracker) SetStatus(id string, s model.DownloadStatus) {
	dt.Lock()
	defer dt.Unlock()
	dt.status[id] = s
}

func (dt *DownloadTracker) SetError(id string, msg string) {
	dt.Lock()
	defer dt.Unlock()
	dt.errs[id] = msg
	dt.status[id] = model.StatusFailed
}

func (dt *DownloadTracker) Get(id string) (model.DownloadStatus, string) {
	dt.Lock()
	defer dt.Unlock()
	return dt.status[id], dt.errs[id]
}

func (dt *DownloadTracker) Delete(id string) {
	dt.Lock()
	defer dt.Unlock()
	delete(dt.status, id)
	delete(dt.errs, id)
}

// Type definition for the main service
type Streamrip struct {
	tracker     *DownloadTracker
	indexer     *Indexer
	fileManager *FileManager
	queries     *db.Queries
}

func NewStreamrip(indexer *Indexer, fileManager *FileManager, queries *db.Queries) *Streamrip {
	return &Streamrip{
		tracker:     NewDownloadTracker(),
		indexer:     indexer,
		fileManager: fileManager,
		queries:     queries,
	}
}

type streamripJSONOutput struct {
	DownloadPath string `json:"downloadPath"`
}

func (s *Streamrip) EnsureTrackForUser(ctx context.Context, songID, user, isrc string, quality int64) (*model.DownloadResult, error) {
	exists, err := s.indexer.IsTrackInLibrary(context.Background(), isrc)
	if err != nil {
		return nil, err
	}
	// Check if already linked
	trackLinkedParams := db.IsTrackLinkedToUserByUsernameAndISRCParams{
		Username: user,
		Isrc:     sql.NullString{String: isrc, Valid: isrc != ""},
	}
	isLinkedInt, err := s.queries.IsTrackLinkedToUserByUsernameAndISRC(context.Background(), trackLinkedParams)
	if err != nil {
		return nil, err
	}
	isLinked := isLinkedInt == 1

	downloadID := uuid.New().String()

	if exists && isLinked {
		return &model.DownloadResult{ID: downloadID, Action: model.ActionNoop}, nil
	} else if exists { // Just get the track info from the db and make the symlink to the apropiate user
		// LinkTrackToUser creates the record for the symlink in the db
		_, err := s.fileManager.LinkTrackToUser(ctx, isrc, user)
		if err != nil {
			return nil, err
		}

		s.tracker.SetStatus(downloadID, model.StatusSuccess)
		track, err := s.queries.SearchTracksByISRC(context.Background(), sql.NullString{String: isrc, Valid: isrc != ""})
		if err != nil && err != sql.ErrNoRows {
			return nil, fmt.Errorf("error searching for the song by ISRC in the DB: %w", err)
		}

		go s.saveDownloadHistory(ctx, downloadID, user, track.ID, quality, string(model.StatusSuccess), "")
		return &model.DownloadResult{ID: downloadID, Action: model.ActionLinked}, nil
	}

	s.tracker.SetStatus(downloadID, model.StatusDownloading)

	// cmd := exec.Command("srip", "--no-db", "id", "qobuz", "track", songID)
	cmd := exec.Command("rip", "--no-db", "id", "qobuz", "track", songID)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.Printf("Executing command: %v", cmd.Args)
	if err := cmd.Start(); err != nil {
		s.tracker.SetError(downloadID, fmt.Sprintf("start error: %v", err))
		//To make sure context doesn't timeout
		ctx = context.Background()
		go s.saveDownloadHistory(ctx, downloadID, user, 0, 0, "failed", err.Error())
		return nil, err
	}

	time.AfterFunc(15*time.Minute, func() {
		s.tracker.Delete(downloadID)
	})

	ctx = context.Background()
	go func() {
		if err := cmd.Wait(); err != nil {
			errMsg := fmt.Sprintf("rip error: %v\n%s", err, stderr.String())
			s.tracker.SetError(downloadID, errMsg)
			s.saveDownloadHistory(ctx, downloadID, user, 0, 0, "failed", errMsg)
			return
		}

		downloadPath, err := extractDownloadPath(stdout.String())
		if err != nil {
			errMsg := fmt.Sprintf("parse error: %v", err)
			s.tracker.SetError(downloadID, errMsg)
			s.saveDownloadHistory(ctx, downloadID, user, 0, 0, "failed", errMsg)
			return
		}

		s.tracker.SetStatus(downloadID, model.StatusIndexing)

		fileInfo, err := os.Stat(downloadPath)
		if err != nil {
			errMsg := fmt.Sprintf("file not found: %v", err)
			s.tracker.SetError(downloadID, errMsg)
			s.saveDownloadHistory(ctx, downloadID, user, 0, 0, "failed", errMsg)
			return
		}

		trackID, err := s.indexer.IndexFile(ctx, fileInfo, downloadPath, user)
		if err != nil {
			errMsg := fmt.Sprintf("indexing error: %v", err)
			s.tracker.SetError(downloadID, errMsg)
			s.saveDownloadHistory(ctx, downloadID, user, 0, 0, "failed", errMsg)
			return
		}

		_, err = s.fileManager.LinkTrackToUser(ctx, isrc, user)
		if err != nil {
			errMsg := fmt.Sprintf("symlink error: %v", err)
			s.tracker.SetError(downloadID, errMsg)
			s.saveDownloadHistory(ctx, downloadID, user, 0, 0, "failed", errMsg)
			return
		}

		s.tracker.SetStatus(downloadID, model.StatusSuccess)
		s.saveDownloadHistory(ctx, downloadID, user, trackID, quality, "success", "")
	}()

	return &model.DownloadResult{ID: downloadID, Action: model.ActionDownloading}, nil
}

func (s *Streamrip) GetDownloadStatus(id string) (model.DownloadStatus, string) {
	return s.tracker.Get(id)
}

func extractDownloadPath(output string) (string, error) {
	start := strings.Index(output, "---BEGIN JSON---")
	end := strings.Index(output, "---END JSON---")

	if start == -1 || end == -1 || start >= end {
		return "", fmt.Errorf("could not find JSON delimiters")
	}

	jsonRaw := output[start+len("---BEGIN JSON---") : end]
	jsonRaw = strings.TrimSpace(jsonRaw)

	var parsed streamripJSONOutput
	if err := json.Unmarshal([]byte(jsonRaw), &parsed); err != nil {
		return "", fmt.Errorf("error parsing json block: %w", err)
	}

	if parsed.DownloadPath == "" {
		return "", fmt.Errorf("parsed JSON has empty downloadPath")
	}

	return parsed.DownloadPath, nil
}

// SearchSong ejecuta una búsqueda usando streamrip y devuelve los resultados en una estructura Go
func (s *Streamrip) SearchSong(source, mediaType, query string) ([]model.StreamripSearchResult, error) {
	// Verificar que el binario existe
	// _, err := exec.LookPath("srip")
	_, err := exec.LookPath("rip")
	if err != nil {
		return nil, errors.New("the command rip is not available in the PATH")
	}

	// Construir el comando
	// cmd := exec.Command("srip", "search", "--stdout", source, mediaType, query)
	cmd := exec.Command("rip", "--no-db", "search", "--stdout", source, mediaType, query)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error executing search: %v\nstderr: %s", err, stderr.String())
	}

	// Parsear el output JSON
	var results []model.StreamripSearchResult
	decoder := json.NewDecoder(strings.NewReader(stdout.String()))
	// decoder.DisallowUnknownFields()

	if err := decoder.Decode(&results); err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %v", err)
	}

	return results, nil
}

func (s *Streamrip) saveDownloadHistory(
	ctx context.Context,
	downloadID, user string,
	trackID int64,
	quality int64,
	status string,
	errorMsg string,
) {
	ctx = context.Background()
	userData, err := s.queries.GetUserByUsername(ctx, user)
	if err != nil {
		log.Printf("Could not find the user %s: %v", user, err)
		return
	}

	params := db.InsertDownloadHistoryParams{
		ID:      downloadID,
		UserID:  sql.NullInt64{Int64: userData.ID, Valid: userData.ID > 0},
		Status:  sql.NullString{String: status, Valid: status != ""},
		Service: sql.NullString{String: "qobuz", Valid: true},
	}

	if errorMsg != "" {
		params.ErrorMessage = sql.NullString{String: errorMsg, Valid: true}
		params.CompletedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}

	if status == "success" {
		params.CompletedAt = sql.NullTime{Time: time.Now(), Valid: true}
		params.TrackID = sql.NullInt64{Int64: trackID, Valid: true}
		params.Quality = sql.NullInt64{Int64: quality, Valid: true}
	}

	_, err = s.queries.InsertDownloadHistory(ctx, params)
	if err != nil {
		log.Printf("Error guardando historial de descarga: %v", err)
	}
}

func (s *Streamrip) GetDeezerTrackSample(isrc string) (sampleUrl string, err error) {
	url := fmt.Sprintf("https://api.deezer.com/track/isrc:%s", isrc)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("error making request to Deezer: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("deezer API returned status %d", resp.StatusCode)
	}

	var result deezerTrackSampleResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding deezer response: %w", err)
	}
	return result.Preview, err
}
