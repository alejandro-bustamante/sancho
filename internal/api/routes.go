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
	GetDownloadStatus(c *gin.Context)
	GetTrackSample(c *gin.Context)
}
type LibraryHandler interface {
	IndexFolder(c *gin.Context)
	GetTracks(c *gin.Context)
	FindTrackInLibrary(c *gin.Context)
	DeleteTrackFromLibrary(c *gin.Context)
	GetUserTracks(c *gin.Context)
	StreamTrack(c *gin.Context)
	GenerateAlbumThumbnails(c *gin.Context)
	GetThumbnailGenerationStatus(c *gin.Context)
}

type UserHandler interface {
	RegisterUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	AuthenticateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
}

func RegisterRoutes(router *gin.Engine, p ProxyHandler, m MusicHandler, l LibraryHandler, u UserHandler) {
	router.Use(mdw.CORSMiddleware())

	api := router.Group("/api")
	{
		api.GET("/proxy", p.ProxyCORSHandler)

		api.POST("/downloads", m.DownloadSingleTrack)
		api.GET("/search", m.SearchTracksByTitle)
		api.GET("/downloads/:id/status", m.GetDownloadStatus)
		api.GET("/search/:isrc/sample", m.GetTrackSample)

		api.POST("/index", l.IndexFolder)
		api.GET("/tracks", l.GetTracks)
		api.GET("/tracks/search", l.FindTrackInLibrary)
		api.DELETE("/users/:username/tracks/:trackId", l.DeleteTrackFromLibrary)

		api.POST("/library/thumbnails", l.GenerateAlbumThumbnails)
		api.GET("/library/thumbnails/status", l.GetThumbnailGenerationStatus)

		api.GET("/users/:username/tracks", l.GetUserTracks)
		api.GET("/tracks/:trackId/stream", l.StreamTrack)

		api.POST("/users", u.RegisterUser)
		api.DELETE("/users", u.DeleteUser)
		api.POST("/auth", u.AuthenticateUser)
		api.PATCH("/users/:id", u.UpdateUser)
	}
}
