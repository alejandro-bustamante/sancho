package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/alejandro-bustamante/sancho/server/internal/api"
	"github.com/alejandro-bustamante/sancho/server/internal/api/controller"
	db "github.com/alejandro-bustamante/sancho/server/internal/repository"
	"github.com/alejandro-bustamante/sancho/server/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3" // Importante: Driver de SQLite
)

func main() {
	// Load enviroment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file. Error: %v", err)
	}

	port := os.Getenv("HTTP_PORT")
	db_path := os.Getenv("DB_PATH")

	// Initialize the db
	// Just one connection open for the whole app
	// The goroutines implementation handle concurrency
	// for the connection behind the scenes
	// conn, err := sql.Open("sqlite3", "database/dev.sancho?_foreign_keys=on")
	conn, err := sql.Open("sqlite3", db_path)
	if err != nil {
		log.Fatalf("Error abriendo la base de datos: %v", err)
	}
	defer conn.Close()
	if err := conn.Ping(); err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}
	queries := db.New(conn)

	// Initialize services and inject dependencies
	indexerService := service.NewIndexer(queries)
	proxyHandler := controller.NewProxyCORSHandler()
	fileMangerService := service.NewFileManager(queries)
	streamripService := service.NewStreamrip(indexerService, fileMangerService, queries)

	// Initialize handlers and inject dependencies
	downloadHandler := controller.NewMusicHandler(streamripService, indexerService, fileMangerService)
	libraryHandler := controller.NewLibraryHandler(queries, indexerService)
	userHandler := controller.NewUserHandler(queries)

	// Configure the router
	router := gin.Default()
	router.Static("/client", "../client")
	api.RegisterRoutes(router, proxyHandler, downloadHandler, libraryHandler, userHandler)
	log.Printf("Server running on http://localhost:%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Could not initialize the server. Error: %v", err)
	}
}
