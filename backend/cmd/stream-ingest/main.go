package main

import (
	"backend/internal/stream-ingest/grpc"
	"backend/internal/stream-ingest/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize gRPC client for the Stream Processing Service
	grpc.InitStreamProcessingClient()

	router := gin.Default()

	// Handlers for stream events and metadata
	router.POST("/start-stream", handlers.HandleStartStream)
	router.POST("/stop-stream", handlers.HandleStopStream)
	router.POST("/forward-metadata", handlers.HandleForwardMetadata)

	log.Println("Stream Ingest Service running on port 8080")
	router.Run(":8080")
}
