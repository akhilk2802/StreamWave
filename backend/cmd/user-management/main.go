// main.go
package main

import (
	"backend/internal/user-management/db"
	"backend/internal/user-management/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()

	router := gin.Default()

	// Define routes for User Management
	router.POST("/users", handlers.CreateUser)
	router.GET("/users/:userId", handlers.GetUser)
	router.PUT("/users/:userId", handlers.UpdateUser)
	router.DELETE("/users/:userId", handlers.DeleteUser)

	// Start the server on port 8080
	router.Run(":8080")
}
