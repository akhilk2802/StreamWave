package routes

import (
	"backend/internal/api/handlers"
	"backend/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", handlers.RegisterHandler)
		authRoutes.POST("/login", handlers.LoginHandler)
	}

	protectedRoutes := r.Group("/")
	protectedRoutes.Use(middleware.AuthMiddleware())
	{
		// protectedRoutes.GET("/streams", handlers.GetStreamsHandler)
		// protectedRoutes.POST("/streams", handlers.CreateStreamHandler)
	}

}
