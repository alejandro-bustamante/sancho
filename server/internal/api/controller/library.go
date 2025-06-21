package controller

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	db "github.com/alejandro-bustamante/sancho/server/internal/repository"
	"github.com/gin-gonic/gin"
)

type Indexer interface {
	IndexFolder(ctx context.Context, rootDir string) error
}

// The structure of the request the client has to pass us
type LibraryIndexRequest struct {
	Path string `json:"path" binding:"required"`
}

type LibraryHandler struct {
	queries        *db.Queries
	indexerService Indexer
}

func NewLibraryHandler(q *db.Queries, s Indexer) *LibraryHandler {
	return &LibraryHandler{
		queries:        q,
		indexerService: s,
	}
}

func (h *LibraryHandler) IndexFolder(c *gin.Context) {
	var req LibraryIndexRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}
	ctx := c.Request.Context()

	go func() {
		if err := h.indexerService.IndexFolder(ctx, req.Path); err != nil {
			log.Printf("Error indexing the folder. Error: %v", err)
		} else {
			log.Printf("Indexing completed for: %s", req.Path)
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"status":  "Indexing in progress",
		"message": "The process is being executed on the background.",
	})
}

func (h *LibraryHandler) GetTracks(c *gin.Context) {
	tracks, err := h.queries.ListTracksByDate(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while getting the tracks"})
		return
	}
	c.JSON(http.StatusOK, tracks)
}

func (h *LibraryHandler) FindTrackInLibrary(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The parameter 'q' was missing"})
		return
	}

	param := sql.NullString{String: query, Valid: true}
	results, err := h.queries.SearchTracksByTitle(c.Request.Context(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while searching the tracks"})
		return
	}
	c.JSON(http.StatusOK, results)
}
