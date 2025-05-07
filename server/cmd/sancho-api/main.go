package main

import (
	"log"

	"github.com/alejandro-bustamante/sancho/server/internal/api"
	"github.com/alejandro-bustamante/sancho/server/internal/api/controller"
	"github.com/alejandro-bustamante/sancho/server/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	//api port
	port := "8081"

	router := gin.Default()
	proxyHandler := controller.NewProxyCORSHandler()
	streamRipService := service.NewStreamripService()
	downloadHandler := controller.NewDownloadHandler(streamRipService)

	router.Static("/client", "../client")

	api.RegisterRoutes(router, proxyHandler, downloadHandler)

	log.Printf("Server running on http://localhost:%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Could not run server: %v", err)
	}
}
