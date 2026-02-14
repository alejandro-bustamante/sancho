package service

import (
	"context"
	"database/sql"

	"github.com/alejandro-bustamante/sancho/server/internal/repository"
)

type LibraryService struct {
	db repository.Database
}

func NewLibraryService(db repository.Database) *LibraryService {
	return &LibraryService{db: db}
}

func (s *LibraryService) GetAllTracks(ctx context.Context) ([]repository.Track, error) {
	return s.db.ListTracksByDate(ctx)
}

func (s *LibraryService) GetUserTracks(ctx context.Context, username string) ([]repository.ListTracksByUsernameRow, error) {
	tracks, err := s.db.ListTracksByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if tracks == nil {
		return []repository.ListTracksByUsernameRow{}, nil
	}
	return tracks, nil
}

func (s *LibraryService) GetTrackByID(ctx context.Context, trackID int64) (repository.Track, error) {
	return s.db.GetTrackByID(ctx, trackID)
}

func (s *LibraryService) SearchTracks(ctx context.Context, query string) ([]repository.Track, error) {
	param := sql.NullString{String: query, Valid: true}
	return s.db.SearchTracksByTitle(ctx, param)
}
