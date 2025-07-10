package config

import (
	"os"
)

func isDev() bool {
	return os.Getenv("SANCHO_ENV") == "dev"
}

var (
	DBPath       string
	SanchoPath   string
	HttpPort     string
	FrontendPath string
)

func init() {
	if isDev() {
		DBPath = os.Getenv("DB_PATH")
		SanchoPath = os.Getenv("SANCHO_PATH")
		HttpPort = os.Getenv("HTTP_PORT")
		FrontendPath = os.Getenv("FRONTEND_PATH")
	} else {
		DBPath = "/data/database.sancho"
		SanchoPath = "/sancho"
		HttpPort = "5400"
		FrontendPath = "/app/build"
	}
}
