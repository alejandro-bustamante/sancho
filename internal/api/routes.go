package api

import (
	"net/http"

	"github.com/alejandro-bustamante/sancho/server/internal/api/controller"
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

func RegisterRoutes(
	mux *http.ServeMux,
	p *controller.ProxyHandler,
	m *controller.MusicHandler,
	l *controller.LibraryHandler,
	u *controller.UserHandler,
) {

	// music handler
	mux.HandleFunc("POST /api/library/index", m.DownloadSingleTrack)
	mux.HandleFunc("GET /api/search", m.SearchTracksByTitle)
	mux.HandleFunc("GET /api/downloads/{id}/status", m.GetDownloadStatus)
	mux.HandleFunc("GET /api/search/{isrc}/sample", m.GetTrackSample)

	// library handler
	mux.HandleFunc("POST /api/index", l.IndexFolder)
	mux.HandleFunc("GET /api/library/thumbnails", l.GenerateAlbumThumbnails)
	mux.HandleFunc("GET /api/tracks", l.GetTracks)
	mux.HandleFunc("GET /api/tracks/search", l.FindTrackInLibrary)
	mux.HandleFunc("GET /api/library/thumbnails/status", l.GetThumbnailGenerationStatus)
	mux.HandleFunc("GET /api/users/{username}/tracks", l.GetUserTracks)
	mux.HandleFunc("GET /api/tracks/{trackId}/stream", l.StreamTrack)
	mux.HandleFunc("DELETE /api/users/{username}/tracks/{trackId}", l.DeleteTrackFromLibrary)

	// user handler
	mux.HandleFunc("GET /", u.HandleIndex)

	mux.HandleFunc("POST /register", u.RegisterUser)
	mux.HandleFunc("POST /login", u.AuthenticateUser)
	mux.HandleFunc("POST /logout", u.LogoutUser)
	mux.HandleFunc("DELETE /delete", u.DeleteUser)
	mux.HandleFunc("DELETE /api/users/{id}", u.UpdateUser)

	// api := router.Group("/api")
	// {
	// 	api.GET("/proxy", p.ProxyCORSHandler)

	// 	api.POST("/downloads", m.DownloadSingleTrack)
	// 	api.GET("/search", m.SearchTracksByTitle)
	// 	api.GET("/downloads/:id/status", m.GetDownloadStatus)
	// 	api.GET("/search/:isrc/sample", m.GetTrackSample)

	// 	api.POST("/index", l.IndexFolder)
	// 	api.POST("/library/thumbnails", l.GenerateAlbumThumbnails)
	// 	api.GET("/tracks", l.GetTracks)
	// 	api.GET("/tracks/search", l.FindTrackInLibrary)
	// 	api.GET("/library/thumbnails/status", l.GetThumbnailGenerationStatus)
	// 	api.GET("/users/:username/tracks", l.GetUserTracks)
	// 	api.GET("/tracks/:trackId/stream", l.StreamTrack)
	// 	api.DELETE("/users/:username/tracks/:trackId", l.DeleteTrackFromLibrary)

	// 	api.POST("/users", u.RegisterUser)
	// 	api.POST("/auth", u.AuthenticateUser)
	// 	api.DELETE("/users", u.DeleteUser)
	// 	api.PATCH("/users/:id", u.UpdateUser)
	// }
}
