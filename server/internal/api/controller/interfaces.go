package controller

import (
	"context"
	"os"

	"github.com/alejandro-bustamante/sancho/server/internal/model"
)

type Streamrip interface {
	EnsureTrackForUser(ctx context.Context, songID, user, isrc string, quality int64) (*model.DownloadResult, error)
	SearchSong(source, mediaType, query string) ([]model.StreamripSearchResult, error)
	GetDownloadStatus(downloadID string) (model.DownloadStatus, string)
}

type Indexer interface {
	IndexFolder(ctx context.Context, rootDir, user string) error
	IndexFile(ctx context.Context, info os.FileInfo, path, user string) (trackID int64, err error)
}

type FileManager interface {
	LinkTrackToUser(ctx context.Context, isrc, user string) (symlinkPath string, err error)
}
