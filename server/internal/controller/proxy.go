package controller

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ProxyDeezerHandler(c *gin.Context) {
	targetURL := c.Query("url")
	if targetURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'url' query parameter"})
		return
	}

	resp, err := http.Get(targetURL)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch external URL"})
		return
	}
	defer resp.Body.Close()

	// Copiamos encabezados relevantes
	c.Status(resp.StatusCode)
	c.Header("Content-Type", resp.Header.Get("Content-Type"))

	// Copiamos el cuerpo
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy response body"})
	}
}
