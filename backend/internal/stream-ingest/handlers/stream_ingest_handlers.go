package handlers

import (
	"backend/internal/stream-ingest/grpc"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleStartStream handles the on_publish event from NGINX
func HandleStartStream(c *gin.Context) {
	streamKey := c.PostForm("name")
	resolution := "1080p"
	format := "dash"

	log.Printf("Stream started with key: %s", streamKey)
	grpc.StartTranscoding(streamKey, resolution, format)

	c.String(http.StatusOK, "Stream processing started")
}

// HandleStopStream handles the on_done event from NGINX
func HandleStopStream(c *gin.Context) {
	streamKey := c.PostForm("name")
	log.Printf("Stream stopped with key: %s", streamKey)

	c.String(http.StatusOK, "Stream stopped")
}

// HandleForwardMetadata forwards metadata to the Stream Processing Service
func HandleForwardMetadata(c *gin.Context) {
	streamKey := c.PostForm("stream_key")
	metadata := c.PostForm("metadata")

	log.Printf("Forwarding metadata for stream: %s", streamKey)
	grpc.ForwardMetadata(streamKey, metadata)

	c.String(http.StatusOK, "Metadata forwarded")
}
