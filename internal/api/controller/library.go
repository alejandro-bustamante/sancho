package controller

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
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

func (h *LibraryHandler) IndexFolder(w http.ResponseWriter, r *http.Request) {
	var req LibraryIndexRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error en el formato de la solicitud", http.StatusBadRequest)
		return
	}

	// Capturamos el path y usuario por si el contexto original se cancela
	path := req.Path
	user := req.User
	service := req.Service
	quality := req.Quality

	// Procesamos en segundo plano
	go func() {
		// Usamos contexto vacío para que no se cancele si el cliente desconecta
		ctx := context.Background()

		log.Printf("Indexing folder '%s' for user '%s'...", path, user)
		// if err := h.indexerService.IndexFolder(ctx, path, user, service, quality); err != nil {
		if err := h.indexerService.IndexFolder(ctx, path, user, service, quality); err != nil {
			log.Printf("[ERROR] Failed indexing folder %s for user %s: %v", path, user, err)
		} else {
			log.Printf("[OK] Indexing completed for folder %s (user: %s)", path, user)
		}
	}()

	// Devolvemos un fragmento HTML para que HTMX lo muestre en pantalla
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("<div class='bg-green-100 p-4'>Indexación iniciada para " + req.User + "</div>"))
}

func (h *LibraryHandler) GetTracks(w http.ResponseWriter, r *http.Request) {
	// tracks, err := h.libraryService.GetAllTracks(r.Context())
	_, err := h.libraryService.GetAllTracks(r.Context())
	if err != nil {
		http.Error(w, "<div class='error'>Error al obtener las canciones</div>", http.StatusInternalServerError)
		return
	}

	// Placeholder HTML. Luego inyectarás un componente Templ pasándole la variable `tracks`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("<ul><li>Placeholder: Se encontraron canciones</li></ul>"))
}

func (h *LibraryHandler) GetUserTracks(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")

	// tracks, err := h.libraryService.GetUserTracks(r.Context(), username)
	_, err := h.libraryService.GetUserTracks(r.Context(), username)
	if err != nil {
		http.Error(w, "<div class='error'>Error obteniendo canciones del usuario</div>", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("<div>Placeholder: Canciones cargadas para " + username + "</div>"))
}

func (h *LibraryHandler) StreamTrack(w http.ResponseWriter, r *http.Request) {
	trackIDStr := r.PathValue("trackId")
	trackID, err := strconv.ParseInt(trackIDStr, 10, 64)
	if err != nil {
		http.Error(w, "ID de canción inválido", http.StatusBadRequest)
		return
	}

	track, err := h.libraryService.GetTrackByID(r.Context(), trackID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Canción no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al obtener canción", http.StatusInternalServerError)
		return
	}

	if _, err := os.Stat(track.FilePath); os.IsNotExist(err) {
		log.Printf("Archivo no encontrado para track %d: %s", trackID, track.FilePath)
		http.Error(w, "Archivo de audio no encontrado", http.StatusNotFound)
		return
	}

	// La librería estándar sirve archivos muy fácilmente
	http.ServeFile(w, r, track.FilePath)
}

func (h *LibraryHandler) FindTrackInLibrary(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "<div class='error'>Falta el parámetro 'q'</div>", http.StatusBadRequest)
		return
	}

	// results, err := h.libraryService.SearchTracks(r.Context(), query)
	_, err := h.libraryService.SearchTracks(r.Context(), query)
	if err != nil {
		http.Error(w, "<div class='error'>Error buscando canciones</div>", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("<div>Placeholder: Resultados de búsqueda en librería</div>"))
}

func (h *LibraryHandler) DeleteTrackFromLibrary(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	trackIDStr := r.PathValue("trackId")

	trackID, err := strconv.ParseInt(trackIDStr, 10, 64)
	if err != nil {
		http.Error(w, "ID de canción inválido", http.StatusBadRequest)
		return
	}

	err = h.fileManager.DeleteTrackForUser(context.Background(), username, trackID)
	if err != nil {
		log.Printf("Error al eliminar la canción: %v", err)
		http.Error(w, "<div class='error'>No se pudo eliminar la canción</div>", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("<div class='success'>Canción eliminada correctamente</div>"))
}

func (h *LibraryHandler) GenerateAlbumThumbnails(w http.ResponseWriter, r *http.Request) {
	h.thumbnailService.GenerateAlbumThumbnails()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("<div>Generación de miniaturas iniciada...</div>"))
}

func (h *LibraryHandler) GetThumbnailGenerationStatus(w http.ResponseWriter, r *http.Request) {
	// isRunning, processed, total, errMsg := h.thumbnailService.GetStatus()
	isRunning, _, _, _ := h.thumbnailService.GetStatus()

	// HTMX puede hacer polling a este endpoint para actualizar una barra de progreso
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("<div>Placeholder Estado: Corriendo=" + strconv.FormatBool(isRunning) + "</div>"))
}
