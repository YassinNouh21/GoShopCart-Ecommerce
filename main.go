package main

import (
	routers "ecommerce/routes"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		// Handle error loading .env file
		return
	}

	// Get the port from the environment variable
	port := os.Getenv("PORT")

	// Create a new Gin router with default middleware
	router := gin.Default()
	router.Use(gin.Logger())

	// Initialize database client and collections
	authRoutes := router.Group("/auth")
	routers.GetAuthRoutes(authRoutes)

	// Set up authentication routes under /auth

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Route not defined",
		})
	})

	// Run the server on the specified port
	router.Run(port)
}
