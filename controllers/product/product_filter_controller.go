package product

import (
	"context"
	"ecommerce/database"
	"ecommerce/models/product"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	errFailedFetchProducts  = errors.New("Failed to fetch products")
	errFailedDecodeProducts = errors.New("Failed to decode products")
	errInvalidMinPrice      = errors.New("Invalid minPrice value")
	errInvalidMaxPrice      = errors.New("Invalid maxPrice value")
	errInvalidPriceRange    = errors.New("Invalid price range")
	errNoProductsFound      = errors.New("No products found")
	errNoPriceProvided      = errors.New("No price provided")
)

// GetProductsByKeyword retrieves products based on a keyword search
func GetProductsByKeyword(c *gin.Context) {
	keyword := c.Query("keyword")

	filter := bson.M{"product_name": bson.M{"$regex": keyword, "$options": "i"}}

	products, err := database.DB.ProductCollection.Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errFailedFetchProducts.Error()})
		c.Abort()
		return
	}

	var result []product.Product
	if err := products.All(context.Background(), &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errFailedDecodeProducts.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetProductsByPriceRange retrieves products within a price range
func GetProductsByPriceRangeController(c *gin.Context) {
	minPriceStr := c.Query("minPrice")
	maxPriceStr := c.Query("maxPrice")
	// Set a timeout for the function execution
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel() // Cancel the context to release resources

	minPrice, err := strconv.ParseFloat(minPriceStr, 32)
	if err != nil && minPriceStr != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidMinPrice.Error()})
		c.Abort()
		return
	}

	maxPrice, err := strconv.ParseFloat(maxPriceStr, 32)

	if err != nil && maxPriceStr != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidMaxPrice.Error()})
		c.Abort()

		return
	}
	if minPriceStr == "" && maxPriceStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidPriceRange.Error()})
		c.Abort()
		return
	}
	var filter primitive.M
	if minPriceStr == "" {
		filter = bson.M{"price": bson.M{"$lte": maxPrice}}

	} else if maxPriceStr == "" {
		filter = bson.M{"price": bson.M{"$gte": minPrice}}
	} else {
		if minPrice > maxPrice {
			c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidPriceRange})
			c.Abort()
			return
		}
		filter = bson.M{"price": bson.M{"$gte": minPrice, "$lte": maxPrice}}
	}

	options := options.Find().SetSort(bson.D{{Key: "price", Value: 1}})

	products, err := database.DB.ProductCollection.Find(ctx, filter, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errFailedFetchProducts.Error()})
		c.Abort()

		return
	}

	var result []product.Product
	if err := products.All(ctx, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errFailedDecodeProducts.Error()})
		c.Abort()

		return
	}

	if len(result) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": errNoProductsFound.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetProductsByPriceRange retrieves products within a price range
func GetProductsByPriceController(c *gin.Context) {
	priceStr := c.Param("price")
	// Set a timeout for the function execution
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel() // Cancel the context to release resources
	if priceStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errNoPriceProvided.Error()})
		c.Abort()
		return
	}
	price, err := strconv.ParseFloat(priceStr, 16)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidPriceRange.Error()})
		c.Abort()
		return
	}

	filter := bson.M{"price": price}
	options := options.Find().SetSort(bson.D{{Key: "price", Value: 1}})

	products, err := database.DB.ProductCollection.Find(ctx, filter, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errFailedFetchProducts.Error()})
		c.Abort()
		return
	}

	var result []product.Product
	if err := products.All(ctx, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errFailedDecodeProducts.Error()})
		c.Abort()
		return
	}

	if len(result) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": errNoProductsFound.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, result)
}
