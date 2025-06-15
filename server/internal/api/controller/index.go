package controller

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	db "github.com/alejandro-bustamante/sancho/server/internal/repository"
	"github.com/alejandro-bustamante/sancho/server/internal/service"
	"github.com/gin-gonic/gin"
)

type IndexRequest struct {
	Path string `json:"path" binding:"required"`
}

type IndexHandler struct{}

func NewIndexHandler() *IndexHandler {
	return &IndexHandler{}
}

func (h *IndexHandler) IndexFolder(c *gin.Context) {
	var req IndexRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	ctx := context.Background()
	dbPath := "database/dev.sancho"

	go func() {
		// if err := service.IndexFolder(ctx, dbPath, req.Path); err != nil {
		if err := service.IndexFolder(ctx, dbPath, req.Path); err != nil {
			log.Printf("Error indexando carpeta: %v", err)
		} else {
			log.Printf("Indexación completada para: %s", req.Path)
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"status":  "Indexación en curso",
		"message": "El proceso se ejecuta en segundo plano.",
	})
}

func (h *IndexHandler) GetTracks(c *gin.Context) {
	dbPath := "database/dev.sancho"
	conn, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo conectar a la base de datos"})
		return
	}
	defer conn.Close()

	queries, err := db.Prepare(context.Background(), conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron preparar las queries"})
		return
	}
	defer queries.Close()

	tracks, err := queries.ListTracksByDate(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los tracks"})
		return
	}

	c.JSON(http.StatusOK, tracks)
}

func (h *IndexHandler) SearchTracks(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Falta el parámetro 'q'"})
		return
	}

	dbPath := "database/dev.sancho"
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo conectar a la base de datos"})
		return
	}
	defer conn.Close()

	// ⚠️ Intentar cargar colación ICU
	// if _, err := conn.Exec(`SELECT icu_load_collation('es_BO', 'french_ai_ci', 'PRIMARY');`); err != nil {
	// 	log.Printf("⚠️ Advertencia: no se pudo cargar la colación ICU: %v", err)
	// }

	queries, err := db.Prepare(context.Background(), conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron preparar las queries"})
		return
	}
	defer queries.Close()

	param := sql.NullString{String: query, Valid: true}
	results, err := queries.SearchTracksByTitle(context.Background(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar tracks"})
		return
	}

	c.JSON(http.StatusOK, results)
}
