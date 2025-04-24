package api

import (
	"github.com/alejandro-bustamante/sancho/server/internal/controller"
	"github.com/alejandro-bustamante/sancho/server/pkg/util"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.Use(util.CORSMiddleware())

	router.GET("/proxy", controller.ProxyDeezerHandler)
}
