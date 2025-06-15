package api

import (
	mdw "github.com/alejandro-bustamante/sancho/server/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

type ProxyHandler interface {
	ProxyCORSHandler(c *gin.Context)
}
type DownloadHandler interface {
	SingleDownloadHandler(c *gin.Context)
	// add more functions from the handler (download.go) when needed
}
type IndexHandler interface {
	IndexFolder(c *gin.Context)
	GetTracks(c *gin.Context)
	SearchTracks(c *gin.Context)
}

func RegisterRoutes(router *gin.Engine, p ProxyHandler, d DownloadHandler, i IndexHandler) {
	router.Use(mdw.CORSMiddleware())

	router.GET("/proxy", p.ProxyCORSHandler)
	router.POST("/download", d.SingleDownloadHandler)
	router.POST("/index", i.IndexFolder)
	router.GET("/tracks", i.GetTracks)

	router.GET("/tracks/search", i.SearchTracks)
}
