package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/alejandro-bustamante/sancho/server/internal/api"
	"github.com/alejandro-bustamante/sancho/server/internal/api/controller"
	"github.com/alejandro-bustamante/sancho/server/internal/api/middleware"
	"github.com/alejandro-bustamante/sancho/server/internal/config"
	db "github.com/alejandro-bustamante/sancho/server/internal/repository"
	"github.com/alejandro-bustamante/sancho/server/internal/service"
	"github.com/spf13/afero"

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
	queries := db.NewDatabase(conn)

	//Creating the afero fs instance for the fileManager and indexer
	fs := afero.NewOsFs()

	// Inicializar servicios
	fileMangerService := service.NewFileManager(queries, fs)
	indexerService := service.NewIndexer(queries, fileMangerService, fs, nil)
	libraryService := service.NewLibraryService(queries)
	streamripService := service.NewStreamrip(indexerService, fileMangerService, queries)
	thumbnailService := service.NewThumbnailService(queries)
	userService := service.NewUserService(queries)

	proxyHandler := controller.NewProxyCORSHandler()

	// Inicializar handlers
	downloadHandler := controller.NewMusicHandler(streamripService, indexerService, fileMangerService)
	libraryHandler := controller.NewLibraryHandler(libraryService, indexerService, fileMangerService, thumbnailService)
	userHandler := controller.NewUserHandler(userService)

	// Configurar router
	// router := gin.Default()
	// Serve library files for album art
	// router.Static("/library", config.LibraryPath)

	// Chage to std lib
	mux := http.NewServeMux()
	mux.Handle("GET /library/", http.StripPrefix("/library/", http.FileServer(http.Dir(config.LibraryPath))))

	// Aquí registramos TODAS las rutas, incluyendo GET / que ahora maneja el UserHandler
	api.RegisterRoutes(mux, proxyHandler, downloadHandler, libraryHandler, userHandler)
	// ------------------------------------------

	// ------------- FRONTEND (ESTÁTICOS) -------------------
	frontend := config.FrontendPath

	// SVELTE CLEANUP: En lugar de /_app, servimos una carpeta general de assets
	// Asegúrate de crear una carpeta 'assets' dentro de tu FrontendPath para tu CSS de Tailwind y JS custom
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(filepath.Join(frontend, "assets")))))

	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(frontend, "favicon.ico"))
	})

	// ELIMINADO: El mux.HandleFunc("GET /", ...) con el placeholder ya no va aquí.
	// Ahora vive en routes.go y lo maneja u.HandleIndex
	// ------------------------------------------

	handlerConCORS := middleware.CORSMiddleware(mux)

	port := config.HttpPort
	log.Printf("Server running on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, handlerConCORS); err != nil {
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
