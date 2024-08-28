package main

import (
	"backend/internal/api/routes"
	"backend/internal/config"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	config.LoadConfig()

	r := gin.Default()

	routes.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("could not run the server: %v", err)
	}

}
