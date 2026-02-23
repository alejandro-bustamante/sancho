package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/alejandro-bustamante/sancho/server/internal/model"
)

type MusicHandler struct {
	streamripService   Streamrip
	indexerService     Indexer
	fileManagerService FileManager
}

func NewMusicHandler(s Streamrip, x Indexer, f FileManager) *MusicHandler {
	return &MusicHandler{
		streamripService:   s,
		indexerService:     x,
		fileManagerService: f,
	}
}

type DownloadRequest struct {
	// Qobuz song's ids in their json are strings
	ID      string `json:"id" binding:"required"`
	ISRC    string `json:"isrc" binding:"required"`
	User    string `json:"user" binding:"required"`
	Quality int64  `json:"quality" binding:"required"`
}

type SearchRequest struct {
	Service   string `json:"service" binding:"required"`
	MediaType string `json:"media_type" binding:"required"`
	Query     string `json:"query" binding:"required"`
}

type TrackSampleRequest struct {
	ISRC string `json:"isrc" binding:"required"`
}

func (h *MusicHandler) DownloadSingleTrack(w http.ResponseWriter, r *http.Request) {
	var req DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "<div class='error'>Solicitud inválida</div>", http.StatusBadRequest)
		return
	}

	log.Printf("Download started for song with Qobuz ID: %s, ISRC: %s", req.ID, req.ISRC)
	result, err := h.streamripService.EnsureTrackForUser(r.Context(), req.ID, req.User, req.ISRC, req.Quality)
	if err != nil {
		log.Printf("Error downloading and indexing song: %v", err)
		http.Error(w, "<div class='error'>Error al iniciar la descarga</div>", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	switch result.Action {
	case model.ActionNoop:
		w.Write([]byte("<div class='info'>La canción ya está en tu cuenta.</div>"))
	case model.ActionLinked:
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("<div class='info'>La canción ya estaba descargada. Vinculada a tu cuenta.</div>"))
	case model.ActionDownloading:
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("<div class='success'>Descarga iniciada.</div>"))
	default:
		http.Error(w, "<div class='error'>Acción desconocida devuelta por el servidor.</div>", http.StatusInternalServerError)
	}
}

func (h *MusicHandler) SearchTracksByTitle(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "<div class='error'>Se requiere el parámetro 'q'</div>", http.StatusBadRequest)
		return
	}

	// results, err := h.streamripService.SearchSong("qobuz", "track", query)
	_, err := h.streamripService.SearchSong("qobuz", "track", query)
	if err != nil {
		http.Error(w, "<div class='error'>Error en la búsqueda</div>", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("<ul><li>Placeholder: Resultados de búsqueda en Qobuz</li></ul>"))
}

func (h *MusicHandler) GetTrackSample(w http.ResponseWriter, r *http.Request) {
	isrc := r.PathValue("isrc")
	sample_url, err := h.streamripService.GetDeezerTrackSample(isrc)
	if err != nil {
		http.Error(w, "Error al obtener sample", http.StatusBadRequest)
		return
	}

	// Para un sample, podrías devolver un tag de audio
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("<audio controls src='" + sample_url + "'></audio>"))
}

func (h *MusicHandler) GetDownloadStatus(w http.ResponseWriter, r *http.Request) {
	downloadID := r.PathValue("id")
	status, _ := h.streamripService.GetDownloadStatus(downloadID)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("<div>Estado de descarga: " + string(status) + "</div>"))
}
