package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	model "github.com/alejandro-bustamante/sancho/server/internal/model"
	db "github.com/alejandro-bustamante/sancho/server/internal/repository"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	tag "go.senan.xyz/taglib"
)

type Indexer struct {
	queries     *db.Queries
	fileManager *FileManager
}

func NewIndexer(queries *db.Queries, fileManager *FileManager) *Indexer {
	return &Indexer{
		queries:     queries,
		fileManager: fileManager,
	}
}

type DeezerIDs struct {
	ArtistID int `json:"artist_id"`
	AlbumID  int `json:"album_id"`
}

func (x *Indexer) IndexFolder(ctx context.Context, rootDir, user, service string, quality int) error {
	ctx = context.Background()
	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !x.isAudioFile(path) {
			return nil
		}

		if err := x.RegisterLocalTrack(ctx, path, user, service, quality); err != nil {
			log.Printf("Error registrando %s: %v", path, err)
		} else {
			log.Printf("✓ Registrado %s", path)
		}

		return nil
	})
}

func (x *Indexer) IndexFile(ctx context.Context, info os.FileInfo, path, user string) (trackID int64, err error) {
	// Get track info
	// Tags: big hash map
	ctx = context.Background()
	tags, err := tag.ReadTags(path)
	if err != nil {
		return 0, fmt.Errorf("error reading tags: %w", err)
	}
	get := func(key string) string {
		vals, ok := tags[key]
		if !ok || len(vals) == 0 {
			return ""
		}
		return vals[0]
	}
	// Properties: some technical data
	properties, err := tag.ReadProperties(path)
	if err != nil {
		return 0, fmt.Errorf("Error reading properties: %w", err)
	}

	// Find deezer ID for artist and album
	isrc := get(tag.ISRC)
	if isrc == "" {
		return 0, fmt.Errorf("could not find the isrc of the provided track: %w", err)
	}
	deezerIDs, err := x.getDeezerIDs(isrc)
	if err != nil {
		return 0, fmt.Errorf("error fetching artist and album id from deezer: %w", err)
	}

	// Verify is the song is not downloaded already
	// This is checked by streamrip service on download, but here for external songs
	existsInt, err := x.queries.TrackExistsByISRC(ctx, sql.NullString{String: isrc, Valid: isrc != ""})
	if err != nil {
		return 0, fmt.Errorf("error checking for the song in the database: %w", err)
	}
	trackExists := existsInt == 1
	if !trackExists {
		// The return value of -1 in the ID is only to indicate this specific failure
		// Yes, this value probably should be sent in the error code. This is faster for now
		return -1, fmt.Errorf("could index the provided track as it already exists in the db: %w", err)
	}

	// Check if we already have artist and album for that song
	artistInDb, err := x.isArtistInDB(ctx, strconv.Itoa(deezerIDs.ArtistID))
	if err != nil {
		return 0, fmt.Errorf("error searching for artist deezer id in the db: %w", err)
	}
	albumInDb, err := x.isAlbumInDB(ctx, strconv.Itoa(deezerIDs.AlbumID))
	if err != nil {
		return 0, fmt.Errorf("error searching for album deezer id in the db: %w", err)
	}

	// Should never happen. Check first
	if albumInDb && !artistInDb {
		return 0, fmt.Errorf("error in db consistency: found album with deezer ID %v with no artist", deezerIDs.AlbumID)
	}

	// We need this var to store the artistID we generate when inserting the artist
	var artistID int64
	// If there is no artist, we insert both artist and album
	if !artistInDb {
		normalizedArtistName, err := NormalizeText(get(tag.Artist))
		if err != nil {
			return 0, fmt.Errorf("error normalizing artists name: %w", err)
		}
		artistParams := db.InsertArtistParams{
			DeezerID:       sql.NullString{String: strconv.Itoa(deezerIDs.ArtistID), Valid: true},
			Name:           get(tag.Artist),
			NormalizedName: normalizedArtistName,
		}
		artist, err := x.queries.InsertArtist(ctx, artistParams)
		if err != nil {
			return 0, fmt.Errorf("error inserting artist into the db %w", err)
		}
		artistID = artist.ID
	} else { // If there is artist we search it and store its ID
		deeezerArtistID := sql.NullString{String: strconv.Itoa(deezerIDs.ArtistID), Valid: true}
		artist, err := x.queries.GetArtistByDeezerID(ctx, deeezerArtistID)
		if err != nil {
			return 0, fmt.Errorf("error searching artist in the db %w", err)
		}
		artistID = artist.ID
	}

	// Same as artistID
	var albumID int64
	// We've determined there IS an artist, we just insert the album
	if !albumInDb {
		//AlbumArtist tags seems to have more compatibility
		normalizedAlbumName, err := NormalizeText(get(tag.Album))
		if err != nil {
			return 0, fmt.Errorf("error normalizing albums name: %w", err)
		}
		totalTracks, err := x.getAlbumTrackNumber(deezerIDs.AlbumID)
		if err != nil {
			return 0, fmt.Errorf("error normalizing albums name: %w", err)
		}

		releaseDate := get(tag.Date)
		genre := get(tag.Genre)
		albumParams := db.InsertAlbumParams{
			DeezerID:        sql.NullString{String: strconv.Itoa(deezerIDs.AlbumID), Valid: true},
			Title:           get(tag.Album),
			NormalizedTitle: normalizedAlbumName,
			ArtistID:        artistID,
			ReleaseDate:     sql.NullString{String: releaseDate, Valid: releaseDate != ""},
			AlbumArtPath:    sql.NullString{String: "", Valid: false},
			Genre:           sql.NullString{String: genre, Valid: genre != ""},
			TotalTracks:     sql.NullInt64{Int64: int64(totalTracks), Valid: totalTracks > 0},
		}
		album, err := x.queries.InsertAlbum(ctx, albumParams)
		if err != nil {
			return 0, fmt.Errorf("error inserting album into the db %w", err)
		}
		albumID = album.ID
	} else {
		deeezerAlbumID := sql.NullString{String: strconv.Itoa(deezerIDs.AlbumID), Valid: true}
		album, err := x.queries.GetAlbumByDeezerID(ctx, deeezerAlbumID)
		if err != nil {
			return 0, fmt.Errorf("error searching artist in the db %w", err)
		}
		albumID = album.ID
	}

	// Insert track
	title := get(tag.Title)
	trackNum := x.parseInt(get(tag.TrackNumber))
	discNum := x.parseInt(get(tag.DiscNumber))
	composer := get(tag.Composer)
	duration := properties.Length.Milliseconds()
	normalizedTrackName, err := NormalizeText(title)
	if err != nil {
		return 0, fmt.Errorf("error normalizing albums name: %w", err)
	}

	params := db.InsertTrackParams{
		Title:           title,
		NormalizedTitle: normalizedTrackName,
		ArtistID:        sql.NullInt64{Int64: artistID, Valid: artistID > 0},
		AlbumID:         sql.NullInt64{Int64: albumID, Valid: albumID > 0},
		Duration:        sql.NullInt64{Int64: duration, Valid: duration > 0},
		TrackNumber:     sql.NullInt64{Int64: trackNum, Valid: trackNum > 0},
		DiscNumber:      sql.NullInt64{Int64: discNum, Valid: discNum > 0},
		SampleRate:      sql.NullInt64{Int64: int64(properties.SampleRate), Valid: properties.SampleRate > 0},
		Bitrate:         sql.NullInt64{Int64: int64(properties.Bitrate), Valid: properties.Bitrate > 0},
		Channels:        sql.NullInt64{Int64: int64(properties.Channels), Valid: properties.Channels > 0},
		FilePath:        path,
		FileSize:        sql.NullInt64{Int64: info.Size(), Valid: true},
		Isrc:            sql.NullString{String: isrc, Valid: isrc != ""},
		Composer:        sql.NullString{String: composer, Valid: composer != ""},
	}
	track, err := x.queries.InsertTrack(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("error storing track: %w", err)
	}

	fmt.Printf("✓ %s (%s)\n", track.Title, track.FilePath)
	return track.ID, nil
}

func (x *Indexer) isArtistInDB(ctx context.Context, deezerID string) (bool, error) {
	existsInt, err := x.queries.ArtistExistsByDeezerID(ctx, sql.NullString{String: deezerID, Valid: true})
	if err != nil {
		return false, fmt.Errorf("db error checking artist existence: %w", err)
	}
	return existsInt == 1, nil
}

func (x *Indexer) isAlbumInDB(ctx context.Context, deezerID string) (bool, error) {
	existsInt, err := x.queries.AlbumExistsByDeezerID(ctx, sql.NullString{String: deezerID, Valid: true})
	if err != nil {
		return false, fmt.Errorf("db error checking album existence: %w", err)
	}
	return existsInt == 1, nil
}

type deezerIDResponse struct {
	Artist struct {
		ID int `json:"id"`
	} `json:"artist"`
	Album struct {
		ID int `json:"id"`
	} `json:"album"`
}

func (x *Indexer) getDeezerIDs(isrc string) (DeezerIDs, error) {
	url := fmt.Sprintf("https://api.deezer.com/track/isrc:%s", isrc)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return DeezerIDs{}, fmt.Errorf("error making request to Deezer: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return DeezerIDs{}, fmt.Errorf("deezer API returned status %d", resp.StatusCode)
	}

	var result deezerIDResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return DeezerIDs{}, fmt.Errorf("error decoding deezer response: %w", err)
	}

	return DeezerIDs{
		ArtistID: result.Artist.ID,
		AlbumID:  result.Album.ID,
	}, nil
}

type deezerNumTracksResponse struct {
	ID int `json:"nb_tracks"`
}

func (x *Indexer) getAlbumTrackNumber(deezerID int) (int, error) {
	url := fmt.Sprintf("https://api.deezer.com/album/%d", deezerID)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("error making request to Deezer: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("deezer API returned status %d", resp.StatusCode)
	}

	var result deezerNumTracksResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("error decoding deezer response: %w", err)
	}
	return result.ID, nil
}

func (x *Indexer) parseInt(s string) int64 {
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}
	return 0
}
func (x *Indexer) isAudioFile(path string) bool {
	ext := filepath.Ext(path)
	switch ext {
	case ".mp3", ".flac", ".ogg", ".m4a", ".wav", ".aiff", ".alac":
		return true
	default:
		return false
	}
}

func (x *Indexer) IsTrackInLibrary(ctx context.Context, isrc string) (bool, error) {
	ctx = context.Background()

	existsInt, err := x.queries.TrackExistsByISRC(ctx, sql.NullString{String: isrc, Valid: isrc != ""})
	if err != nil {
		return false, fmt.Errorf("db error checking track existence: %w", err)

	}
	return existsInt == 1, nil
}

func (x *Indexer) RegisterLocalTrack(ctx context.Context, fullPath, user, service string, quality int) error {
	ctx = context.Background()
	info, err := os.Stat(fullPath)
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	trackID, err := x.IndexFile(ctx, info, fullPath, user)
	if trackID == -1 {
		// This id is returned in case we find the songs id already int the db
		return fmt.Errorf("song found already in the db, skiping...: %w", err)
	}
	if err != nil {
		return fmt.Errorf("indexing error: %w", err)
	}

	track, err := x.queries.GetTrackByID(ctx, trackID)
	if err != nil {
		return fmt.Errorf("could not fetch track after indexing: %w", err)
	}

	artist, err := x.queries.GetArtistByTrackID(ctx, trackID)
	if err != nil {
		return fmt.Errorf("error fetching artist: %w", err)
	}

	trackModel := model.TrackFromDB(track)
	renamedPath, err := x.fileManager.renameTrack(ctx, trackModel, artist.Name)
	if err != nil {
		return fmt.Errorf("error renaming track: %w", err)
	}
	trackModel.FilePath = renamedPath

	_, err = x.fileManager.moveTrackToLibrary(ctx, trackModel)
	if err != nil {
		return fmt.Errorf("error moving track to library: %w", err)
	}

	_, err = x.fileManager.LinkTrackToUser(ctx, *trackModel.ISRC, user)
	if err != nil {
		return fmt.Errorf("error linking track to user: %w", err)
	}

	x.saveTransferHistory(ctx, user, service, trackID, quality)

	return nil
}

func (x *Indexer) saveTransferHistory(
	ctx context.Context,
	user, service string,
	trackID int64,
	quality int,
) {
	ctx = context.Background()
	userData, err := x.queries.GetUserByUsername(ctx, user)
	if err != nil {
		log.Printf("Could not find the user %s: %v", user, err)
		return
	}
	downloadID := uuid.New().String()
	status := string(model.StatusTransfered)

	params := db.InsertDownloadHistoryParams{
		ID:          downloadID,
		UserID:      sql.NullInt64{Int64: userData.ID, Valid: userData.ID > 0},
		TrackID:     sql.NullInt64{Int64: trackID, Valid: trackID > 0},
		Quality:     sql.NullInt64{Int64: int64(quality), Valid: quality >= 0},
		Status:      sql.NullString{String: status, Valid: status != ""},
		Service:     sql.NullString{String: service, Valid: service != ""},
		CompletedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	}

	_, err = x.queries.InsertDownloadHistory(ctx, params)
	if err != nil {
		log.Printf("Error guardando historial de descarga: %v", err)
	}
}
