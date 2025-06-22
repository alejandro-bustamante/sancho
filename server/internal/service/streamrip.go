package service

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"time"

	"encoding/json"
	"errors"
	"strings"

	model "github.com/alejandro-bustamante/sancho/server/internal/model"
)

type StreamripProd struct{}

func NewStreamripService() *StreamripProd {
	return &StreamripProd{}
}

func (s *StreamripProd) DownloadTrack(url, title, artist, album, user string) (string, error) {
	// Generar un ID único para la descarga
	downloadID := fmt.Sprintf("%d", time.Now().UnixNano())

	// Verificar que el comando rip existe
	_, err := exec.LookPath("srip")
	if err != nil {
		return "", fmt.Errorf("sancho-streamrip (srip command) is not installed or not present on the PATH: %w", err)
	}

	// Configurar path de salida basado en el usuario
	// outputPath := fmt.Sprintf("Downloads/%s", user)

	// Capturar salida del comando para debug
	var stdout, stderr bytes.Buffer

	// Construir comando de streamrip CORRECTAMENTE
	// Notar que 'url' es un subcomando, seguido de la URL real
	cmd := exec.Command(
		"rip",
		"url", // Subcomando para URLs
		url,   // La URL real como argumento
	)

	// Configurar captura de salida
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Lanzar comando
	log.Printf("Executin command: %v", cmd.Args)
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("error starting download: %w", err)
	}

	// Iniciar una goroutine para capturar el resultado final sin bloquear la respuesta HTTP
	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Printf("Error in the download %s: %v\nStderr: %s", downloadID, err, stderr.String())
		} else {
			log.Printf("Download %s completed successfully \nOutput: %s", downloadID, stdout.String())
		}
	}()

	return downloadID, nil
}

// SearchSong ejecuta una búsqueda usando streamrip y devuelve los resultados en una estructura Go
func (s *StreamripProd) SearchSong(source, mediaType, query string) ([]model.StreamripSearchResult, error) {
	// Verificar que el binario existe
	_, err := exec.LookPath("srip")
	if err != nil {
		return nil, errors.New("the command srip is not available in the PATH")
	}

	// Construir el comando
	cmd := exec.Command("srip", "search", "--stdout", source, mediaType, query)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error executing search: %v\nstderr: %s", err, stderr.String())
	}

	// Parsear el output JSON
	var results []model.StreamripSearchResult
	decoder := json.NewDecoder(strings.NewReader(stdout.String()))
	// decoder.DisallowUnknownFields()

	if err := decoder.Decode(&results); err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %v", err)
	}

	return results, nil
}
