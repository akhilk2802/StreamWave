package routes

import (
	"backend/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/ping", handlers.PingHandler)

	// Define other routes here
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/login", handlers.LoginHandler)
		authRoutes.POST("/register", handlers.RegisterHandler)
	}

	streamRoutes := r.Group("/streams")
	{
		streamRoutes.GET("/", handlers.GetStreamsHandler)
		streamRoutes.POST("/", handlers.CreateStreamHandler)
	}
}
