package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/alejandro-bustamante/sancho/server/internal/db"
	_ "github.com/mattn/go-sqlite3"

	taglib "go.senan.xyz/taglib"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Uso: %s <directorio_musica>", os.Args[0])
	}
	rootDir := os.Args[1]

	// Abrir base de datos
	conn, err := sql.Open("sqlite3", "database/dev.sancho")
	if err != nil {
		log.Fatalf("error abriendo base de datos: %v", err)
	}
	defer conn.Close()

	ctx := context.Background()
	queries, err := db.Prepare(ctx, conn)
	if err != nil {
		log.Fatalf("error preparando queries: %v", err)
	}
	defer queries.Close()

	// Recorrer archivos
	err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if !isAudioFile(path) {
			return nil
		}
		if err := indexFile(ctx, path, info, queries); err != nil {
			log.Printf("Error procesando %s: %v", path, err)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error recorriendo directorio: %v", err)
	}
}

func indexFile(ctx context.Context, path string, info os.FileInfo, queries *db.Queries) error {
	tags, err := taglib.ReadTags(path)
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

	// Parseos
	title := get(taglib.Title)
	if title == "" {
		title = filepath.Base(path)
	}
	trackNum := parseInt(get(taglib.TrackNumber))
	discNum := parseInt(get(taglib.DiscNumber))
	duration := parseInt(get(taglib.Length)) // en milisegundos

	params := db.InsertTrackParams{
		Title:       title,
		FilePath:    path,
		FileSize:    sql.NullInt64{Int64: info.Size(), Valid: true},
		Codec:       sql.NullString{String: get(taglib.Encoding), Valid: get(taglib.Encoding) != ""},
		Duration:    sql.NullInt64{Int64: duration, Valid: duration > 0},
		Bitrate:     sql.NullInt64{}, // No disponible directamente
		Channels:    sql.NullInt64{},
		SampleRate:  sql.NullInt64{},
		BitDepth:    sql.NullInt64{},
		TrackNumber: sql.NullInt64{Int64: trackNum, Valid: trackNum > 0},
		DiscNumber:  sql.NullInt64{Int64: discNum, Valid: discNum > 0},
		ArtistID:    sql.NullInt64{}, // Resolver con búsqueda
		AlbumID:     sql.NullInt64{}, // Resolver con búsqueda
		Isrc:        sql.NullString{String: get(taglib.ISRC), Valid: get(taglib.ISRC) != ""},
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
