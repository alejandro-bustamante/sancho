package controller

import (
	"context"
	"os"

	"github.com/alejandro-bustamante/sancho/server/internal/model"
)

type Streamrip interface {
	DownloadTrack(ctx context.Context, songID, user string, quality int64) (string, error)
	SearchSong(source, mediaType, query string) ([]model.StreamripSearchResult, error)
	GetDownloadStatus(downloadID string) (model.DownloadStatus, string)
}

type Indexer interface {
	IndexFolder(ctx context.Context, rootDir, user string) error
	IndexFile(ctx context.Context, info os.FileInfo, path, user string) (trackID int64, err error)
}
