package config

import (
	"os"
)

func isDev() bool {
	return os.Getenv("SANCHO_ENV") == "dev"
}

var (
	DBPath     string
	SanchoPath string
	HttpPort   string
)

func init() {
	if isDev() {
		DBPath = os.Getenv("DB_PATH")
		SanchoPath = os.Getenv("SANCHO_PATH")
		HttpPort = os.Getenv("HTTP_PORT")
	} else {
		DBPath = "/data/database.sancho"
		SanchoPath = "/sancho"
		HttpPort = "8081"
	}
}
