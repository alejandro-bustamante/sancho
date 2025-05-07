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
}

func RegisterRoutes(router *gin.Engine, p ProxyHandler, d DownloadHandler) {
	router.Use(mdw.CORSMiddleware())

	router.GET("/proxy", p.ProxyCORSHandler)
	router.POST("/download", d.SingleDownloadHandler)
}
