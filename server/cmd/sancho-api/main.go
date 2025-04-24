package main

import (
	"log"

	"github.com/alejandro-bustamante/sancho/server/api"
	"github.com/gin-gonic/gin"
)

func main() {
	//api port
	port := "8081"

	router := gin.Default()

	api.RegisterRoutes(router)

	log.Printf("Server running on http://localhost:%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Could not run server: %v", err)
	}
}
