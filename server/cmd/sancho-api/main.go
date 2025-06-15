package main

import (
	"database/sql"
	"log"

	"github.com/alejandro-bustamante/sancho/server/internal/api"
	"github.com/alejandro-bustamante/sancho/server/internal/api/controller"
	db "github.com/alejandro-bustamante/sancho/server/internal/repository"
	"github.com/alejandro-bustamante/sancho/server/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	conn, err := sql.Open("sqlite3", "database/dev.sancho")
	if err != nil {
		log.Fatal("Error opening the database. Error: %i", err)
	}
	defer conn.Close()

	queries := db.New(conn)
	defer queries.Close()

	//api port
	port := "8081"

	router := gin.Default()
	proxyHandler := controller.NewProxyCORSHandler()
	streamRipService := service.NewStreamripService()
	downloadHandler := controller.NewDownloadHandler(streamRipService)
	indexHandler := controller.NewIndexHandler()

	router.Static("/client", "../client")

	api.RegisterRoutes(router, proxyHandler, downloadHandler, indexHandler)

	log.Printf("Server running on http://localhost:%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Could not run server: %v", err)
	}
}
