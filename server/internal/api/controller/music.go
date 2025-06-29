package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/alejandro-bustamante/sancho/server/internal/model"
	"github.com/gin-gonic/gin"
)

type MusicHandler struct {
	streamripService Streamrip
	indexerService   Indexer
}

func NewMusicHandler(s Streamrip, x Indexer) *MusicHandler {
	return &MusicHandler{
		streamripService: s,
		indexerService:   x,
	}
}

type DownloadRequest struct {
	// Qobuz song's ids in their json are strings
	ID      string `json:"id" binding:"required"`
	ISRC    string `json:"isrc" binding:"required"`
	User    string `json:"user" binding:"required"`
	Quality int64  `json:"quality" binding:"required"`
}
type SearchTrackRequest struct {
	Title string `json:"title" binding:"required"`
}
type SearchRequest struct {
	Service   string `json:"service" binding:"required"`
	MediaType string `json:"media_type" binding:"required"`
	Query     string `json:"query" binding:"required"`
}

func (h *MusicHandler) DownloadSingleTrack(c *gin.Context) {
	var req DownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	log.Printf("Download started for song with Qobuz ID: %s, ISRC: %s", req.ID, req.ISRC)
	downloadID, err := h.streamripService.DownloadTrack(c.Request.Context(), req.ID, req.User, req.Quality)
	if err != nil {
		log.Printf("Error downloading and indexing song: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start download", "details": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"downloadId": downloadID,
		"status":     "downloading",
		"message":    "Download has started, it's state can be checken on with the downloadID.",
	})
}

func (h *MusicHandler) SearchTracksByTitle(c *gin.Context) {
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

func (h *MusicHandler) Search(c *gin.Context) {
	var req SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	results, err := h.streamripService.SearchSong(req.Service, req.MediaType, req.Query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The search failed", "details": err.Error()})
		return
	}
	switch req.MediaType {
	case "track":
		preview := model.MapToTrackPreviews(results)
		c.JSON(http.StatusOK, gin.H{
			"message": "Seach completed succesfully",
			"results": preview,
		})
	default:
		return
	}
}

func (h *MusicHandler) SearchTracksDeezer(c *gin.Context) {
	var req SearchTrackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	query := url.QueryEscape(req.Title)
	deezerURL := "https://api.deezer.com/search?q=" + query

	resp, err := http.Get(deezerURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch from Deezer", "details": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Deezer API returned non-200", "status": resp.Status})
		return
	}

	var deezerResp model.DeezerSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&deezerResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse Deezer response", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Resultados de Deezer",
		"results": deezerResp.Data,
	})
}

func (h *MusicHandler) GetDownloadStatus(c *gin.Context) {
	downloadID := c.Param("id")
	status, errMsg := h.streamripService.GetDownloadStatus(downloadID)

	resp := gin.H{"downloadId": downloadID, "status": status}
	if status == model.StatusFailed {
		resp["error"] = errMsg
	}
	c.JSON(http.StatusOK, resp)
}
