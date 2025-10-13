package service

import (
	"context"
	"database/sql"
	"log"
	"os/exec"
	"path/filepath"
	"sync"

	db "github.com/alejandro-bustamante/sancho/server/internal/repository"
)

// ThumbnailGenerationTracker holds the state of the generation process.
type ThumbnailGenerationTracker struct {
	sync.Mutex
	IsRunning bool
	Total     int
	Processed int
	Error     string
}

type ThumbnailService struct {
	queries *db.Queries
	tracker *ThumbnailGenerationTracker
}

func NewThumbnailService(queries *db.Queries) *ThumbnailService {
	return &ThumbnailService{
		queries: queries,
		tracker: &ThumbnailGenerationTracker{},
	}
}

func (s *ThumbnailService) GetStatus() (bool, int, int, string) {
	s.tracker.Lock()
	defer s.tracker.Unlock()
	return s.tracker.IsRunning, s.tracker.Processed, s.tracker.Total, s.tracker.Error
}

func (s *ThumbnailService) GenerateAlbumThumbnails() {
	s.tracker.Lock()
	if s.tracker.IsRunning {
		s.tracker.Unlock()
		log.Println("Thumbnail generation is already in progress.")
		return
	}
	s.tracker.IsRunning = true
	s.tracker.Processed = 0
	s.tracker.Total = 0
	s.tracker.Error = ""
	s.tracker.Unlock()

	go func() {
		defer func() {
			s.tracker.Lock()
			s.tracker.IsRunning = false
			s.tracker.Unlock()
		}()

		ctx := context.Background()

		albums, err := s.queries.GetAlbumsWithoutArt(ctx)
		if err != nil {
			log.Printf("Error fetching albums without art: %v", err)
			s.tracker.Lock()
			s.tracker.Error = "Failed to fetch albums from database."
			s.tracker.Unlock()
			return
		}

		s.tracker.Lock()
		s.tracker.Total = len(albums)
		s.tracker.Unlock()

		for _, album := range albums {
			track, err := s.queries.GetFirstTrackByAlbumID(ctx, sql.NullInt64{Int64: album.ID, Valid: true})
			if err != nil {
				log.Printf("Could not find a track for album ID %d (%s). Marking as processed. Error: %v", album.ID, album.Title, err)
				s.updateAlbumArtPath(ctx, album.ID, "/dev/null")
			} else {
				albumDir := filepath.Dir(track.FilePath)
				thumbnailPath := filepath.Join(albumDir, "cover_thumb.jpg")

				cmd := exec.Command("ffmpeg", "-hide_banner", "-loglevel", "error", "-i", track.FilePath, "-map", "0:v:0", "-frames:v", "1", "-vf", "scale=iw/4:ih/4", "-update", "1", "-y", thumbnailPath)

				output, err := cmd.CombinedOutput()
				if err != nil {
					log.Printf("Failed to generate thumbnail for album ID %d (%s). FFMPEG error: %s. Error: %v", album.ID, album.Title, string(output), err)
					s.updateAlbumArtPath(ctx, album.ID, "/dev/null")
				} else {
					log.Printf("Generated thumbnail for album '%s' at %s", album.Title, thumbnailPath)
					s.updateAlbumArtPath(ctx, album.ID, thumbnailPath)
				}
			}

			s.tracker.Lock()
			s.tracker.Processed++
			s.tracker.Unlock()
		}
		log.Println("Finished generating album thumbnails.")
	}()
}

func (s *ThumbnailService) updateAlbumArtPath(ctx context.Context, albumID int64, path string) {
	updateParams := db.UpdateAlbumArtPathParams{
		ID:           albumID,
		AlbumArtPath: sql.NullString{String: path, Valid: true},
	}
	if err := s.queries.UpdateAlbumArtPath(ctx, updateParams); err != nil {
		log.Printf("Failed to update album art path for album ID %d: %v", albumID, err)
	}
}
