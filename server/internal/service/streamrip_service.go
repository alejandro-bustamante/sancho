package service

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"time"
)

type StreamripService interface {
	DownloadTrack(url, title, artist, album, user string) (string, error)
}

type StreamripServiceProd struct{}

func NewStreamripService() *StreamripServiceProd {
	return &StreamripServiceProd{}
}

func (s *StreamripServiceProd) DownloadTrack(url, title, artist, album, user string) (string, error) {
	// Generar un ID único para la descarga
	downloadID := fmt.Sprintf("%d", time.Now().UnixNano())
	
	// Verificar que el comando rip existe
	_, err := exec.LookPath("rip")
	if err != nil {
		return "", fmt.Errorf("streamrip (rip command) no está instalado o no está en el PATH: %w", err)
	}
	
	// Configurar path de salida basado en el usuario
	// outputPath := fmt.Sprintf("Downloads/%s", user)
	
	// Capturar salida del comando para debug
	var stdout, stderr bytes.Buffer
	
	// Construir comando de streamrip CORRECTAMENTE
	// Notar que 'url' es un subcomando, seguido de la URL real
	cmd := exec.Command(
		"rip",
		"url",        // Subcomando para URLs
		url,          // La URL real como argumento
		// "--output", outputPath,
	)
	
	// Configurar captura de salida
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	// Lanzar comando
	log.Printf("Ejecutando comando: %v", cmd.Args)
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("error iniciando descarga: %w", err)
	}
	
	// Iniciar una goroutine para capturar el resultado final sin bloquear la respuesta HTTP
	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Printf("Error en descarga %s: %v\nStderr: %s", downloadID, err, stderr.String())
		} else {
			log.Printf("Descarga %s completada exitosamente\nSalida: %s", downloadID, stdout.String())
		}
	}()
	
	return downloadID, nil
}