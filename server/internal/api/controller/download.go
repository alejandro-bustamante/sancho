package controller

import (
	"log"
	"net/http"

	"github.com/alejandro-bustamante/sancho/server/internal/service"
	"github.com/gin-gonic/gin"
)

type DownloadHandlerProd struct {
	streamripService service.StreamripService
}

func NewDownloadHandler(s service.StreamripService) *DownloadHandlerProd {
	return &DownloadHandlerProd{
		streamripService: s,
	}
}

type DownloadRequest struct {
	URL    string `json:"url" binding:"required"`
	Title  string `json:"title" binding:"required"`
	Artist string `json:"artist" binding:"required"`
	Album  string `json:"album" binding:"required"`
	User   string `json:"user" binding:"required"`
}

func (h *DownloadHandlerProd) SingleDownloadHandler(c *gin.Context) {
	var req DownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}
	
	log.Printf("Iniciando descarga para URL: %s, Título: %s, Artista: %s", req.URL, req.Title, req.Artist)
	
	downloadID, err := h.streamripService.DownloadTrack(req.URL, req.Title, req.Artist, req.Album, req.User)
	if err != nil {
		log.Printf("Error al iniciar descarga: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start download", "details": err.Error()})
		return
	}
	
	log.Printf("Descarga iniciada con ID: %s", downloadID)
	
	// Responder con JSON para proporcionar más información
	c.JSON(http.StatusOK, gin.H{
		"downloadId": downloadID,
		"status": "downloading",
		"message": "Descarga iniciada correctamente",
		"trackInfo": gin.H{
			"title": req.Title,
			"artist": req.Artist,
			"album": req.Album,
		},
	})
}