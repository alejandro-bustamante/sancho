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
	SearchTracksByTitle(c *gin.Context)
	SearchTracksDeezer(c *gin.Context)
}
type LibraryHandler interface {
	IndexFolder(c *gin.Context)
	GetTracks(c *gin.Context)
	FindTrackInLibrary(c *gin.Context)
}

func RegisterRoutes(router *gin.Engine, p ProxyHandler, d DownloadHandler, l LibraryHandler) {
	// Cors middleware
	router.Use(mdw.CORSMiddleware())
	router.GET("/proxy", p.ProxyCORSHandler)

	// Download enpoints (works with streamrip)
	router.POST("/download", d.DownloadSingleTrack)
	router.POST("/search", d.SearchTracksByTitle)
	router.POST("/search/deezer", d.SearchTracksDeezer)

	// Library endpoints (works with internal db)
	router.POST("/index", l.IndexFolder)
	router.GET("/tracks", l.GetTracks)
	router.GET("/tracks/search", l.FindTrackInLibrary)

}
