package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	db "github.com/alejandro-bustamante/sancho/server/internal/repository"
	_ "github.com/mattn/go-sqlite3"

	tag "go.senan.xyz/taglib"
)

func IndexFolder(ctx context.Context, dbPath string, rootDir string) error {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("error abriendo base de datos: %w", err)
	}
	defer conn.Close()

	queries, err := db.Prepare(ctx, conn)
	if err != nil {
		return fmt.Errorf("error preparando queries: %w", err)
	}
	defer queries.Close()

	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !isAudioFile(path) {
			return nil
		}
		if err := indexFile(ctx, path, info, queries); err != nil {
			log.Printf("Error procesando %s: %v", path, err)
		}
		return nil
	})
}

func indexFile(ctx context.Context, path string, info os.FileInfo, queries *db.Queries) error {
	tags, err := tag.ReadTags(path)
	if err != nil {
		return fmt.Errorf("error leyendo tags: %w", err)
	}

	get := func(key string) string {
		vals, ok := tags[key]
		if !ok || len(vals) == 0 {
			return ""
		}
		return vals[0]
	}

	title := get(tag.Title)
	if title == "" {
		title = filepath.Base(path)
	}
	trackNum := parseInt(get(tag.TrackNumber))
	discNum := parseInt(get(tag.DiscNumber))
	duration := parseInt(get(tag.Length))

	properties, err := tag.ReadProperties(path)
	if err != nil {
		return fmt.Errorf("error leyendo propiedades del archivo: %w", err)
	}

	params := db.InsertTrackParams{
		Title:       title,
		FilePath:    path,
		FileSize:    sql.NullInt64{Int64: info.Size(), Valid: true},
		Codec:       sql.NullString{String: get(tag.Album), Valid: get(tag.Album) != ""},
		Duration:    sql.NullInt64{Int64: duration, Valid: duration > 0},
		Bitrate:     sql.NullInt64{Int64: int64(properties.Bitrate), Valid: properties.Bitrate > 0},
		Channels:    sql.NullInt64{Int64: int64(properties.Channels), Valid: properties.Channels > 0},
		SampleRate:  sql.NullInt64{Int64: int64(properties.SampleRate), Valid: properties.SampleRate > 0},
		BitDepth:    sql.NullInt64{}, // No disponible
		TrackNumber: sql.NullInt64{Int64: trackNum, Valid: trackNum > 0},
		DiscNumber:  sql.NullInt64{Int64: discNum, Valid: discNum > 0},
		ArtistID:    sql.NullInt64{}, // Resolver luego
		AlbumID:     sql.NullInt64{},
		Isrc:        sql.NullString{String: get(tag.ISRC), Valid: get(tag.ISRC) != ""},
	}

	track, err := queries.InsertTrack(ctx, params)
	if err != nil {
		return fmt.Errorf("error insertando track: %w", err)
	}

	fmt.Printf("✓ %s (%s)\n", track.Title, track.FilePath)
	return nil
}

func parseInt(s string) int64 {
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}
	return 0
}

func isAudioFile(path string) bool {
	ext := filepath.Ext(path)
	switch ext {
	case ".mp3", ".flac", ".ogg", ".m4a", ".wav", ".aiff", ".alac":
		return true
	default:
		return false
	}
}
