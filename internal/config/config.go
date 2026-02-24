package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var (
	DBPath       string
	SanchoPath   string
	HttpPort     string
	FrontendPath string
	LibraryPath  string
)

func init() {
	_ = godotenv.Load()

	DBPath = getEnvOrDefault("DB_PATH", "/data/database.sancho")
	SanchoPath = getEnvOrDefault("SANCHO_PATH", "/sancho")
	HttpPort = getEnvOrDefault("HTTP_PORT", "5400")
	//Probably change frotend path in prod
	FrontendPath = getEnvOrDefault("FRONTEND_PATH", "/app/build")
	LibraryPath = getEnvOrDefault("LIBRARY_PATH", filepath.Join(SanchoPath, "library"))
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
