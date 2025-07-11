package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	tag "go.senan.xyz/taglib"
)

// Estructura para parsear respuesta JSON de Deezer
type DeezerTrack struct {
	Title         string `json:"title"`
	TrackPosition int    `json:"track_position"`
	ReleaseDate   string `json:"release_date"`
	ISRC          string `json:"isrc"`
	Artist        struct {
		Name string `json:"name"`
	} `json:"artist"`
	Album struct {
		Title string `json:"title"`
	} `json:"album"`
}

func getTag(tags map[string][]string, key string) string {
	vals, ok := tags[key]
	if !ok || len(vals) == 0 {
		return ""
	}
	return vals[0]
}

func fetchDeezerData(isrc string) (*DeezerTrack, error) {
	url := fmt.Sprintf("https://api.deezer.com/track/isrc:%s", isrc)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching Deezer data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected response from Deezer: %s", string(body))
	}

	var track DeezerTrack
	if err := json.NewDecoder(resp.Body).Decode(&track); err != nil {
		return nil, fmt.Errorf("error decoding Deezer JSON: %w", err)
	}

	return &track, nil
}

func updateTags(path string) error {
	// Leer tags existentes
	tags, err := tag.ReadTags(path)
	if err != nil {
		return fmt.Errorf("error reading tags: %w", err)
	}

	isrc := getTag(tags, tag.ISRC)
	if isrc == "" {
		return fmt.Errorf("no ISRC found in file: %s", path)
	}

	fmt.Printf("🎧 [%s] ISRC: %s\n", filepath.Base(path), isrc)

	// Espera 0.5 segundos para evitar saturar la API
	time.Sleep(500 * time.Millisecond)

	// Obtener datos desde Deezer
	track, err := fetchDeezerData(isrc)
	if err != nil {
		return fmt.Errorf("deezer lookup failed: %w", err)
	}

	// Crear nuevos tags para sobrescribir
	newTags := map[string][]string{
		tag.Title:       {track.Title},
		tag.Artist:      {track.Artist.Name},
		tag.Album:       {track.Album.Title},
		tag.ISRC:        {track.ISRC},
		tag.TrackNumber: {fmt.Sprintf("%d", track.TrackPosition)},
		tag.Date:        {track.ReleaseDate},
	}

	if err := tag.WriteTags(path, newTags, 0); err != nil {
		return fmt.Errorf("error writing tags: %w", err)
	}

	fmt.Printf("✅ Updated: %s - %s\n", track.Artist.Name, track.Title)
	return nil
}

func processFolder(root string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".flac" || ext == ".mp3" {
			if err := updateTags(path); err != nil {
				fmt.Fprintf(os.Stderr, "⚠️  Error on %s: %v\n", path, err)
			}
		}
		return nil
	})
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <folder>")
		os.Exit(1)
	}

	folder := os.Args[1]
	if err := processFolder(folder); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}
}
