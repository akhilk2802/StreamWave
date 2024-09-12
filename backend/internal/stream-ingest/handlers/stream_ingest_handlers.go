package handlers

import (
	"backend/internal/stream-ingest/services"
	"backend/proto"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleStartStream(c *gin.Context) {
	streamKey := c.PostForm("name") // Extract stream key from POST request
	// app := c.PostForm("app")         // Extract application (optional)

	grpcService := services.StreamService{}

	// Call the gRPC StartStream method
	_, err := grpcService.StartStream(context.Background(), &proto.StreamRequest{
		StreamKey: streamKey,
		// UserId:    "",  // You can add User ID logic here if needed
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failure",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Stream started",
	})

}

func HandleStopStream(c *gin.Context) {
	streamKey := c.PostForm("name") // Extract stream key from POST request

	grpcService := services.StreamService{}

	// Call the gRPC StopStream method
	_, err := grpcService.StopStream(context.Background(), &proto.StreamRequest{
		StreamKey: streamKey,
		UserId:    "", // You can add User ID logic here if needed
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failure",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Stream stopped",
	})
}
