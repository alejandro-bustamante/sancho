package controller

import (
	"context"
	"os"

	"github.com/alejandro-bustamante/sancho/server/internal/model"
	"github.com/alejandro-bustamante/sancho/server/internal/repository"
)

type Streamrip interface {
	EnsureTrackForUser(ctx context.Context, songID, user, isrc string, quality int64) (*model.DownloadResult, error)
	SearchSong(source, mediaType, query string) ([]model.StreamripSearchResult, error)
	GetDownloadStatus(downloadID string) (model.DownloadStatus, string)
	GetDeezerTrackSample(isrc string) (sampleUrl string, err error)
}

type Indexer interface {
	// IndexFolder(ctx context.Context, rootDir, user string) error
	IndexFolder(ctx context.Context, rootDir, user, service string, quality int) error
	IndexFile(ctx context.Context, info os.FileInfo, path, user string) (trackID int64, err error)
}

type FileManager interface {
	LinkTrackToUser(ctx context.Context, isrc, user string) (symlinkPath string, err error)
	DeleteTrackForUser(ctx context.Context, username string, trackID int64) error
}

type ThumbnailService interface {
	GenerateAlbumThumbnails()
	GetStatus() (isRunning bool, processed int, total int, err string)
}

type UserService interface {
	RegisterUser(ctx context.Context, username, password, email string) (*model.User, error)
	AuthenticateUser(ctx context.Context, username, password string) (bool, error)
}

type LibraryService interface {
	GetAllTracks(ctx context.Context) ([]repository.Track, error)
	GetUserTracks(ctx context.Context, username string) ([]repository.ListTracksByUsernameRow, error)
	GetTrackByID(ctx context.Context, trackID int64) (repository.Track, error)
	SearchTracks(ctx context.Context, query string) ([]repository.Track, error)
}
