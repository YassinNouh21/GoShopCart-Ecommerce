package main

import (
	"github.com/YassinNouh21/GoShopCart-Ecommerce/database"
	"github.com/YassinNouh21/GoShopCart-Ecommerce/middlewares"
	routers "github.com/YassinNouh21/GoShopCart-Ecommerce/routes"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func envPortOr(port string) string {
	// If `PORT` variable in environment exists, return it
	if envPort := os.Getenv("PORT"); envPort != "" {
		return ":" + envPort
	}
	// Otherwise, return the value of `port` variable from function argument
	return ":" + port
}

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
	routers.ProfileRoutes(userRoutes)
	routers.AddressRoutes(userRoutes)
	routers.CartRoutes(userRoutes)
	// Set up product-related routes under /product
	productRoutes := router.Group("/product")
	routers.ProductRoutes(productRoutes)
	routers.ProductFilterRoutes(productRoutes)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Route not defined",
		})
	})

	// Run the server on the specified port
	router.Run(":" + port)
}
