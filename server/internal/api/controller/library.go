package controller

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// The structure of the request the client has to pass us
type LibraryIndexRequest struct {
	Path    string `json:"path" binding:"required"`
	User    string `json:"user" binding:"required"`
	Service string `json:"service" binding:"required"`
	Quality int    `json:"quality" binding:"required"`
}

type LibraryHandler struct {
	libraryService   LibraryService
	indexerService   Indexer
	fileManager      FileManager
	thumbnailService ThumbnailService
}

func NewLibraryHandler(lib LibraryService, s Indexer, f FileManager, t ThumbnailService) *LibraryHandler {
	return &LibraryHandler{
		libraryService:   lib,
		indexerService:   s,
		fileManager:      f,
		thumbnailService: t,
	}
}

func (h *LibraryHandler) IndexFolder(c *gin.Context) {
	var req LibraryIndexRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Capturamos el path y usuario por si el contexto original se cancela
	path := req.Path
	user := req.User
	service := req.Service
	quality := req.Quality

	// Respondemos inmediatamente
	c.JSON(http.StatusAccepted, gin.H{
		"status":  "Indexing in progress",
		"message": fmt.Sprintf("Indexing of %s started in background for user %s", path, user),
	})

	// Procesamos en segundo plano
	go func() {
		// Usamos contexto vacío para que no se cancele si el cliente desconecta
		ctx := context.Background()

		log.Printf("Indexing folder '%s' for user '%s'...", path, user)
		if err := h.indexerService.IndexFolder(ctx, path, user, service, quality); err != nil {
			log.Printf("[ERROR] Failed indexing folder %s for user %s: %v", path, user, err)
		} else {
			log.Printf("[OK] Indexing completed for folder %s (user: %s)", path, user)
		}
	}()
}

func (h *LibraryHandler) GetTracks(c *gin.Context) {
	// Uso del servicio
	tracks, err := h.libraryService.GetAllTracks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while getting the tracks"})
		return
	}
	c.JSON(http.StatusOK, tracks)
}

func (h *LibraryHandler) GetUserTracks(c *gin.Context) {
	username := c.Param("username")

	tracks, err := h.libraryService.GetUserTracks(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user tracks"})
		return
	}

	c.JSON(http.StatusOK, tracks)
}

func (h *LibraryHandler) StreamTrack(c *gin.Context) {
	trackIDStr := c.Param("trackId")
	trackID, err := strconv.ParseInt(trackIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid track ID"})
		return
	}

	// Uso del servicio para obtener metadata
	track, err := h.libraryService.GetTrackByID(c.Request.Context(), trackID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "track not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get track"})
		return
	}

	if _, err := os.Stat(track.FilePath); os.IsNotExist(err) {
		log.Printf("File not found for track %d: %s", trackID, track.FilePath)
		c.JSON(http.StatusNotFound, gin.H{"error": "audio file not found"})
		return
	}

	c.File(track.FilePath)
}

func (h *LibraryHandler) FindTrackInLibrary(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The parameter 'q' was missing"})
		return
	}

	results, err := h.libraryService.SearchTracks(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while searching the tracks"})
		return
	}
	c.JSON(http.StatusOK, results)
}

func (h *LibraryHandler) DeleteTrackFromLibrary(c *gin.Context) {
	username := c.Param("username")
	trackIDStr := c.Param("trackId")

	trackID, err := strconv.ParseInt(trackIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": " invalid song ID"})
		return
	}

	ctx := context.Background()
	err = h.fileManager.DeleteTrackForUser(ctx, username, trackID)
	if err != nil {
		log.Printf("Error al eliminar la canción: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar la canción", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Canción eliminada de la librería correctamente"})
}

func (h *LibraryHandler) GenerateAlbumThumbnails(c *gin.Context) {
	h.thumbnailService.GenerateAlbumThumbnails()
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Album thumbnail generation started in the background.",
	})
}

func (h *LibraryHandler) GetThumbnailGenerationStatus(c *gin.Context) {
	isRunning, processed, total, errMsg := h.thumbnailService.GetStatus()
	c.JSON(http.StatusOK, gin.H{
		"isRunning": isRunning,
		"processed": processed,
		"total":     total,
		"error":     errMsg,
	})
}
