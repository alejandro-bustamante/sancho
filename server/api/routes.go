package api

import (
	"github.com/alejandro-bustamante/sancho/server/pkg/util"
	"github.com/gin-gonic/gin"
)

type ProxyHandler interface {
	ProxyCORSHandler(c *gin.Context)
}

func RegisterRoutes(router *gin.Engine, h ProxyHandler) {
	router.Use(util.CORSMiddleware())

	router.GET("/proxy", h.ProxyCORSHandler)
}
