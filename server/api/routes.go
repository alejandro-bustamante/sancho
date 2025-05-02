package api

import (
	"github.com/alejandro-bustamante/sancho/server/pkg/util"
	"github.com/gin-gonic/gin"
)

type ProxyHandler interface {
	ProxyCORSHandler(c *gin.Context)
}
type DownloadHandler interface {
	SingleDownloadHandler(c *gin.Context)
}

func RegisterRoutes(router *gin.Engine, p ProxyHandler, d DownloadHandler) {
	router.Use(util.CORSMiddleware())

	router.GET("/proxy", p.ProxyCORSHandler)
	router.POST("/download", d.SingleDownloadHandler)
}
