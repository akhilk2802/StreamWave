package main

import (
	grpcServer "backend/internal/stream-ingest/grpc"
	"backend/internal/stream-ingest/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Start the gRPC server in a separate goroutine
	go func() {
		grpcServer.StartGRPCServer()
	}()

	// Set up the HTTP server for NGINX RTMP hooks
	r := gin.Default()

	// Route for starting a stream (on_publish hook)
	r.POST("/start-stream", handlers.HandleStartStream)

	// Route for stopping a stream (on_done hook)
	r.POST("/stop-stream", handlers.HandleStopStream)

	// Start the HTTP server on port 8081 (since 8080 is used by your auth service)
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Failed to run HTTP server: %v", err)
	}
}
