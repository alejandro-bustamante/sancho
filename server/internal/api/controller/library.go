package controller

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	db "github.com/alejandro-bustamante/sancho/server/internal/repository"
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
