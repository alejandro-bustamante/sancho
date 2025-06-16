package api

import (
	mdw "github.com/alejandro-bustamante/sancho/server/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

type ProxyHandler interface {
	ProxyCORSHandler(c *gin.Context)
}
type DownloadHandler interface {
	DownloadSingleTrack(c *gin.Context)
	// add more functions from the handler (download.go) when needed
}
type LibraryHandler interface {
	IndexFolder(c *gin.Context)
	GetTracks(c *gin.Context)
	SearchTracks(c *gin.Context)
}

func RegisterRoutes(router *gin.Engine, p ProxyHandler, d DownloadHandler, l LibraryHandler) {
	router.Use(mdw.CORSMiddleware())

	router.GET("/proxy", p.ProxyCORSHandler)
	router.POST("/download", d.DownloadSingleTrack)
	router.POST("/index", l.IndexFolder)

	router.GET("/tracks", l.GetTracks)

	router.GET("/tracks/search", l.SearchTracks)
}
