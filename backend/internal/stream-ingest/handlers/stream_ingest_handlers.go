package handlers

import (
	"log"
	"net/http"

	grpcclient "backend/internal/stream-ingest/grpc"

	"github.com/gin-gonic/gin"
)

// HandleStartStream handles the on_publish hook from NGINX RTMP
func HandleStartStream(c *gin.Context) {
	streamKey := c.PostForm("name")
	log.Printf("Starting stream with key: %s", streamKey)

	// Start the stream processing via gRPC
	if err := grpcclient.StartStreamProcessing(streamKey); err != nil {
		log.Printf("Error starting stream processing: %v", err)
		c.String(http.StatusInternalServerError, "Error starting stream")
		return
	}

	c.String(http.StatusOK, "Stream started successfully")
}

// HandleStopStream handles the on_done hook from NGINX RTMP
func HandleStopStream(c *gin.Context) {
	streamKey := c.PostForm("name")
	log.Printf("Stopping stream with key: %s", streamKey)

	// Implement stop logic if needed

	c.String(http.StatusOK, "Stream stopped successfully")
}
