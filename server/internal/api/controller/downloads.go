package controller

import (
	"log"
	"net/http"

	"github.com/alejandro-bustamante/sancho/server/internal/model"
	"github.com/gin-gonic/gin"
)

type Streamrip interface {
	DownloadTrack(url, title, artist, album, user string) (string, error)
	SearchSong(source, mediaType, query string) ([]model.StreamripSearchResult, error)
}

type DownloadHandlerProd struct {
	streamripService Streamrip
}

func NewDownloadHandler(s Streamrip) *DownloadHandlerProd {
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
type SearchTrackRequest struct {
	Title string `json:"title" binding:"required"`
}

func (h *DownloadHandlerProd) DownloadSingleTrack(c *gin.Context) {
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
		"status":     "downloading",
		"message":    "Descarga iniciada correctamente",
		"trackInfo": gin.H{
			"title":  req.Title,
			"artist": req.Artist,
			"album":  req.Album,
		},
	})
}

func (h *DownloadHandlerProd) SearchTracksByTitle(c *gin.Context) {
	var req SearchTrackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	results, err := h.streamripService.SearchSong("qobuz", "track", req.Title)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to search the track", "details": err.Error()})
		return
	}
	preview := model.MapToTrackPreviews(results)
	c.JSON(http.StatusOK, gin.H{
		"message": "Búsqueda completada",
		"results": preview,
	})

}
