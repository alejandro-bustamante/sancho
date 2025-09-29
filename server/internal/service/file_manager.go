package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/alejandro-bustamante/sancho/server/internal/config"
	model "github.com/alejandro-bustamante/sancho/server/internal/model"
	db "github.com/alejandro-bustamante/sancho/server/internal/repository"
)

type FileManager struct {
	db      *sql.DB
	queries *db.Queries
}

func NewFileManager(db *sql.DB, queries *db.Queries) *FileManager {
	return &FileManager{
		db:      db,
		queries: queries,
	}
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

func (fm *FileManager) DeleteTrackForUser(ctx context.Context, username string, trackID int64) error {
	ctx = context.Background()
	user, err := fm.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("error finding user '%s': %w", username, err)
	}

	// Begin a new database transaction
	tx, err := fm.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create a new querier that operates within the transaction
	// Use qtx for any operation that has be done in this transaction
	qtx := fm.queries.WithTx(tx)

	track, err := qtx.GetTrackByID(ctx, trackID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("track with ID %d does not exist", trackID)
		}
		return fmt.Errorf("error finding track: %w", err)
	}

	// Get the user-track relationship to find the symlink details
	// Check if the song really exists within the user's library
	_, err = qtx.GetUserTrack(ctx, db.GetUserTrackParams{
		UserID:  sql.NullInt64{Int64: user.ID, Valid: true},
		TrackID: sql.NullInt64{Int64: trackID, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("user does not have this track in their library: %w", err)
	}

	// --- Filesystem Operations ---

	// 1. Remove the symlink from the user's personal library folder
	// Find the full path of the symlink
	sanchoRoot := config.SanchoPath
	userLibraryDir := filepath.Join(sanchoRoot, fmt.Sprintf("%s_library", user.Username))
	globalLibraryDir := filepath.Join(sanchoRoot, "library")

	relativeTrackPath, err := filepath.Rel(globalLibraryDir, track.FilePath)
	if err != nil {
		return fmt.Errorf("error calculating relative path: %w", err)
	}
	symlinkPath := filepath.Join(userLibraryDir, relativeTrackPath)

	if err := os.Remove(symlinkPath); err != nil {
		// Log as a warning because the DB state is more critical
		log.Printf("Warning: could not remove symlink %s: %v", symlinkPath, err)
	}

	// --- Database Operations ---

	// 2. Delete the entry from the user_track table
	err = qtx.DeleteUserTrack(ctx, db.DeleteUserTrackParams{
		UserID:  sql.NullInt64{Int64: user.ID, Valid: true},
		TrackID: sql.NullInt64{Int64: trackID, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("error deleting user-track relationship: %w", err)
	}

	// 3. Check if any other users have the same track
	count, err := qtx.CountUsersForTrack(ctx, sql.NullInt64{Int64: trackID, Valid: true})
	if err != nil {
		return fmt.Errorf("error counting users for track: %w", err)
	}

	// 4. If no one else has the track, remove the physical file and database records
	if count == 0 {
		log.Printf("No other user has the track '%s'. Proceeding with deletion from the global library.", track.Title)

		// Remove the physical audio file
		if err := os.Remove(track.FilePath); err != nil {
			log.Printf("Warning: could not remove physical file %s: %v", track.FilePath, err)
		}

		// Delete the track from the 'track' table
		if err := qtx.DeleteTrack(ctx, track.ID); err != nil {
			return fmt.Errorf("error deleting track from database: %w", err)
		}

		// Clean up orphaned album and artist records
		if track.AlbumID.Valid {
			albumTracks, _ := qtx.CountTracksInAlbum(ctx, track.AlbumID)
			if albumTracks == 0 {
				log.Printf("Album ID %d has no tracks left. Deleting...", track.AlbumID.Int64)
				qtx.DeleteAlbum(ctx, track.AlbumID.Int64)
			}
		}
		if track.ArtistID.Valid {
			artistAlbums, _ := qtx.CountAlbumsByArtist(ctx, track.ArtistID.Int64)
			// For a more robust check, you could also count remaining tracks for the artist
			if artistAlbums == 0 {
				log.Printf("Artist ID %d has no albums left. Deleting...", track.ArtistID.Int64)
				qtx.DeleteArtist(ctx, track.ArtistID.Int64)
			}
		}
	}

	// If all operations were successful, commit the transaction
	return tx.Commit()
}
