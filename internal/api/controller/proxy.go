package controller

import (
	"io"
	"net/http"
)

type ProxyHandler struct{}

func NewProxyCORSHandler() *ProxyHandler {
	return &ProxyHandler{}
}

func (h *ProxyHandler) ProxyCORSHandler(w http.ResponseWriter, r *http.Request) {
	targetURL := r.URL.Query().Get("url")
	if targetURL == "" {
		http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
		return
	}

	resp, err := http.Get(targetURL)
	if err != nil {
		http.Error(w, "Failed to fetch external URL", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copiamos encabezados relevantes
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)

	// Copiamos el cuerpo de la imagen directamente al cliente
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		// Nota: Si ya escribiste el WriteHeader, no puedes enviar http.Error aquí
		return
	}
}
