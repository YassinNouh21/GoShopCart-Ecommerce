package main

import (
	"ecommerce/database"
	"ecommerce/middlewares"
	routers "ecommerce/routes"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// initializeDB initializes the database client and collections.
func initializeDB() {
	database.InitializeMongoDBCollections()
}

func main() {
	initializeDB()
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

	// Use Authentication middleware
	router.Use(middlewares.Authentication())

	// Set up user-related routes under /user and product-related routes under /product
	userRoutes := router.Group("/user")
	routers.AddressRoutes(userRoutes)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Route not defined",
		})
	})

	// Run the server on the specified port
	router.Run(port)
}
