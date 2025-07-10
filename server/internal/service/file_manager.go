package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/alejandro-bustamante/sancho/server/internal/config"
	model "github.com/alejandro-bustamante/sancho/server/internal/model"
	db "github.com/alejandro-bustamante/sancho/server/internal/repository"
)

type FileManager struct {
	queries *db.Queries
}

func NewFileManager(queries *db.Queries) *FileManager {
	return &FileManager{queries: queries}
}

func (fm *FileManager) renameTrack(ctx context.Context, track model.Track, artistName string) (string, error) {
	//Temporal empty context to avoid timeouts
	ctx = context.Background()

	baseDir := filepath.Dir(track.FilePath)
	ext := filepath.Ext(track.FilePath)
	// Expected format for the song:
	// <TrackNumber>. <TrackName> - <Artist><extension>
	// Example: 04. Call Me - Blondie.flac

	trackNumber := "00"
	if track.TrackNumber != nil {
		trackNumber = fmt.Sprintf("%02d", *track.TrackNumber)
	}

	safeTitle := sanitizeFilename(track.Title)
	safeArtist := sanitizeFilename(artistName)

	newFileName := fmt.Sprintf("%s. %s - %s%s", trackNumber, safeTitle, safeArtist, ext)
	newPath := filepath.Join(baseDir, newFileName)

	if err := os.Rename(track.FilePath, newPath); err != nil {
		return "", fmt.Errorf("error renaming file: %w", err)
	}

	err := fm.queries.UpdateTrackFilePath(ctx, db.UpdateTrackFilePathParams{
		FilePath: newPath,
		TrackID:  track.ID,
	})
	if err != nil {
		return "", fmt.Errorf("error updating file path in DB: %w", err)
	}

	return newPath, nil
}

func (fm *FileManager) moveTrackToLibrary(ctx context.Context, track model.Track) (trackPath string, err error) {
	//Temporal empty context to avoid timeouts
	ctx = context.Background()
	// Expected folder structure in the library is:
	// Library
	// ├── Artist Name
	// │   ├── Album A
	// │   │   ├── Song 1.flac
	// │   │   ├── Song 2.flac
	// │   └── Album B
	// │       ├── Track 1.m4a
	// └────── └── Track 3.m4a
	sanchoRoot := config.SanchoPath
	libraryRoot := filepath.Join(sanchoRoot, "library")

	artist, err := fm.queries.GetArtistByTrackID(ctx, track.ID)
	if err != nil {
		return "", fmt.Errorf("error fetching artist: %w", err)
	}
	album, err := fm.queries.GetAlbumByTrackID(ctx, track.ID)
	if err != nil {
		return "", fmt.Errorf("error fetching album: %w", err)
	}

	safeArtist := sanitizeFilename(artist.Name)
	safeAlbum := sanitizeFilename(album.Title)

	targetDir := filepath.Join(libraryRoot, safeArtist, safeAlbum)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("error creating directory structure: %w", err)
	}

	fileName := filepath.Base(track.FilePath)
	newPath := filepath.Join(targetDir, fileName)

	if err := os.Rename(track.FilePath, newPath); err != nil {
		return "", fmt.Errorf("error moving file to library: %w", err)
	}

	err = fm.queries.UpdateTrackFilePath(ctx, db.UpdateTrackFilePathParams{
		FilePath: newPath,
		TrackID:  track.ID,
	})
	if err != nil {
		return "", fmt.Errorf("error updating moved file path in DB: %w", err)
	}

	return newPath, nil
}

func (fm *FileManager) LinkTrackToUser(ctx context.Context, isrc, user string) (string, error) {
	//Temporal empty context to avoid timeouts
	ctx = context.Background()
	isrcNull := sql.NullString{String: isrc, Valid: true}

	trackDB, err := fm.queries.SearchTracksByISRC(ctx, isrcNull)
	if err != nil {
		return "", fmt.Errorf("error searching for track: %w", err)
	}

	artist, err := fm.queries.GetArtistByTrackID(ctx, trackDB.ID)
	if err != nil {
		return "", fmt.Errorf("error fetching artist: %w", err)
	}

	trackModel := model.TrackFromDB(trackDB)

	renamedPath, err := fm.renameTrack(ctx, trackModel, artist.Name)
	if err != nil {
		return "", err
	}
	trackModel.FilePath = renamedPath

	finalPath, err := fm.moveTrackToLibrary(ctx, trackModel)
	if err != nil {
		return "", err
	}

	sanchoRoot := config.SanchoPath
	libraryRoot := filepath.Join(sanchoRoot, "library")

	relativeTrackPath, err := filepath.Rel(libraryRoot, finalPath)
	if err != nil {
		return "", fmt.Errorf("error computing relative track path: %w", err)
	}

	userLibraryDir := filepath.Join(sanchoRoot, fmt.Sprintf("%s_library", user))
	userFilePath := filepath.Join(userLibraryDir, relativeTrackPath)
	userDir := filepath.Dir(userFilePath)

	if err := os.MkdirAll(userDir, 0755); err != nil {
		return "", fmt.Errorf("error creating user directory: %w", err)
	}

	// Generate relative symlink path (target path relative to symlink location)
	relativeSymlinkTarget, err := filepath.Rel(userDir, finalPath)
	if err != nil {
		return "", fmt.Errorf("error generating relative symlink target: %w", err)
	}

	if err := os.Symlink(relativeSymlinkTarget, userFilePath); err != nil {
		return "", fmt.Errorf("error creating symlink: %w", err)
	}

	userDB, err := fm.queries.GetUserByUsername(ctx, user)
	if err != nil {
		return "", fmt.Errorf("error searching user in the DB: %w", err)
	}
	trackUserParams := db.AddTrackToUserParams{
		UserID:      sql.NullInt64{Int64: userDB.ID, Valid: userDB.ID > 0},
		TrackID:     sql.NullInt64{Int64: trackDB.ID, Valid: trackDB.ID > 0},
		SymlinkPath: relativeSymlinkTarget, // <-- aquí se guarda el path RELATIVO
	}
	err = fm.queries.AddTrackToUser(ctx, trackUserParams)
	if err != nil {
		return "", fmt.Errorf("error adding row to user_track table: %w", err)
	}

	return userFilePath, nil
}

func sanitizeFilename(name string) string {
	// Chars known to cause issues
	invalidChars := regexp.MustCompile(`[<>:"/\\|?*]`)
	return invalidChars.ReplaceAllString(name, "")
}
