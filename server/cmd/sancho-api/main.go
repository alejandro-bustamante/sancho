package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/alejandro-bustamante/sancho/server/internal/api"
	"github.com/alejandro-bustamante/sancho/server/internal/api/controller"
	"github.com/alejandro-bustamante/sancho/server/internal/config"
	db "github.com/alejandro-bustamante/sancho/server/internal/repository"
	"github.com/alejandro-bustamante/sancho/server/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	dbPath := config.DBPath

	if err := os.MkdirAll(filepath.Dir(dbPath), os.ModePerm); err != nil {
		log.Fatalf("No se pudo crear el directorio para la base de datos: %v", err)
	}
	if err := runMigrations(dbPath); err != nil {
		log.Fatalf("Error ejecutando migraciones: %v", err)
	}

	// ------------- BACKEND -------------------
	// Initialize the db
	// Just one connection open for the whole app
	// The goroutines implementation handle concurrency
	// for the connection behind the scenes

	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Error abriendo la base de datos: %v", err)
	}
	defer conn.Close()
	if err := conn.Ping(); err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}
	queries := db.New(conn)

	// Inicializar servicios
	fileMangerService := service.NewFileManager(queries)
	indexerService := service.NewIndexer(queries, fileMangerService)
	proxyHandler := controller.NewProxyCORSHandler()
	streamripService := service.NewStreamrip(indexerService, fileMangerService, queries)

	// Inicializar handlers
	downloadHandler := controller.NewMusicHandler(streamripService, indexerService, fileMangerService)
	libraryHandler := controller.NewLibraryHandler(queries, indexerService)
	userHandler := controller.NewUserHandler(queries)

	// Configurar router
	router := gin.Default()
	api.RegisterRoutes(router, proxyHandler, downloadHandler, libraryHandler, userHandler)
	// ------------------------------------------

	// ------------- FRONTEND -------------------
	// Servir archivos específicos
	frontend := config.FrontendPath
	// router.Static("/_app", "./build/_app")
	router.Static("/_app", filepath.Join(frontend, "_app"))
	// router.StaticFile("/favicon.png", "./build/favicon.png")
	router.StaticFile("/favicon.png", filepath.Join(frontend, "favicon.png"))

	// Servir el archivo index.html para todas las rutas que no sean API o archivos estáticos (SPA)
	router.NoRoute(func(c *gin.Context) {
		indexPath := filepath.Join(frontend, "index.html")
		c.File(indexPath)
	})
	// ------------------------------------------

	port := config.HttpPort
	log.Printf("Server running on http://localhost:%s", port)
	// if err := router.Run(":" + port); err != nil {
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Could not initialize the server. Error: %v", err)
	}
}

func runMigrations(dbPath string) error {
	dbURL := "file:" + dbPath + "?cache=shared&_fk=1"

	conn, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	driver, err := sqlite.WithInstance(conn, &sqlite.Config{})
	if err != nil {
		return err
	}

	// Asume que tus archivos de migración están en ./migrations
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"sqlite3", driver)
	if err != nil {
		return err
	}

	// Ejecuta todas las migraciones pendientes
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// TODO:
// 1. Standardize the time format in the db. sqlite uses UTC, the backend sends local time
// 3. Polling

// 2. Auth - seems complete
