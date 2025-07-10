package api

import (
	mdw "github.com/alejandro-bustamante/sancho/server/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

type ProxyHandler interface {
	ProxyCORSHandler(c *gin.Context)
}
type MusicHandler interface {
	DownloadSingleTrack(c *gin.Context)
	SearchTracksByTitle(c *gin.Context)
	// SearchTracksDeezer(c *gin.Context)
	GetDownloadStatus(c *gin.Context)
	GetTrackSample(c *gin.Context)
}
type LibraryHandler interface {
	IndexFolder(c *gin.Context)
	GetTracks(c *gin.Context)
	FindTrackInLibrary(c *gin.Context)
}

type UserHandler interface {
	RegisterUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	AuthenticateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
}

func RegisterRoutes(router *gin.Engine, p ProxyHandler, m MusicHandler, l LibraryHandler, u UserHandler) {
	router.Use(mdw.CORSMiddleware())
	router.GET("/proxy", p.ProxyCORSHandler)

	router.POST("/downloads", m.DownloadSingleTrack)
	router.GET("/search", m.SearchTracksByTitle)
	// router.POST("/search/deezer", m.SearchTracksDeezer)
	router.GET("/downloads/:id/status", m.GetDownloadStatus)
	router.GET("/search/:isrc/sample", m.GetTrackSample)

	router.POST("/index", l.IndexFolder)
	router.GET("/tracks", l.GetTracks)
	router.GET("/tracks/search", l.FindTrackInLibrary)

	router.POST("/users", u.RegisterUser)
	router.DELETE("/users", u.DeleteUser)
	router.POST("/auth", u.AuthenticateUser)
	router.PATCH("/users/:id", u.UpdateUser)
}
