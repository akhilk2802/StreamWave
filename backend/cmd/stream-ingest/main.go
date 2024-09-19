package main

import (
	"log"

	"backend/internal/stream-ingest/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
		// Proceeding without .env file
	}

	router := gin.Default()

	// Define routes for NGINX RTMP hooks
	router.POST("/start-stream", handlers.HandleStartStream)
	router.POST("/stop-stream", handlers.HandleStopStream)

	log.Println("Stream Ingest Service running on port 8081")
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Failed to run HTTP server: %v", err)
	}
}
