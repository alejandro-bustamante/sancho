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

// IndexerProd ahora tiene una dependencia de la base de datos (queries).
type IndexerProd struct {
	queries *db.Queries
}

// El constructor ahora requiere que se le pasen las queries.
func NewIndexerService(queries *db.Queries) *IndexerProd {
	return &IndexerProd{
		queries: queries,
	}
}

// IndexFolder ya no necesita el 'dbPath', usa las queries internas.
func (x *IndexerProd) IndexFolder(ctx context.Context, rootDir string) error {
	// ¡La lógica de conexión a la BD ha sido eliminada!
	// Ahora el servicio es más limpio y se enfoca solo en su tarea.
	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !x.isAudioFile(path) {
			return nil
		}
		// Se pasan las queries que ya tiene el struct.
		if err := x.indexFile(ctx, path, info, x.queries); err != nil {
			log.Printf("Error procesando %s: %v", path, err)
		}
		return nil
	})
}

// indexFile no cambia, pero es llamado por un IndexFolder refactorizado.
func (x *IndexerProd) indexFile(ctx context.Context, path string, info os.FileInfo, queries *db.Queries) error {
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
	trackNum := x.parseInt(get(tag.TrackNumber))
	discNum := x.parseInt(get(tag.DiscNumber))
	duration := x.parseInt(get(tag.Length))
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

func (x *IndexerProd) parseInt(s string) int64 {
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}
	return 0
}
func (x *IndexerProd) isAudioFile(path string) bool {
	ext := filepath.Ext(path)
	switch ext {
	case ".mp3", ".flac", ".ogg", ".m4a", ".wav", ".aiff", ".alac":
		return true
	default:
		return false
	}
}
