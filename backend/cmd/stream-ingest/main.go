package main

import (
	"backend/internal/stream-ingest/grpc"
	"backend/internal/stream-ingest/handlers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("STREAM_INGEST_PORT")

	// Initialize gRPC client for the Stream Processing Service
	grpc.InitStreamProcessingClient()

	router := gin.Default()

	// Handlers for stream events and metadata
	router.POST("/start-stream", handlers.HandleStartStream)
	router.POST("/stop-stream", handlers.HandleStopStream)
	router.POST("/forward-metadata", handlers.HandleForwardMetadata)

	log.Printf("Stream Ingest Service running on port %s", port)
	router.Run(":" + port)
}
