package product

import (
	"context"
	"ecommerce/database"
	productModel "ecommerce/models/product"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
CreateProductController handles the creation of a new product.

	It binds the request body to the Product model and returns an error if the request body is invalid.
	It checks if the product already exists and returns an error if it does.
	It creates the new product and returns the ID of the inserted product.
*/
func CreateProductController(c *gin.Context) {
	// Bind the request body to the Product model

	var product productModel.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		c.Abort()
		return
	}
	if product.ProductID != primitive.NilObjectID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot provide productID in the request body"})
		c.Abort()
		return
	}
	product.ProductID = primitive.NewObjectID()

	// Check if the product already exists
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	existingProduct := productModel.Product{}
	err := database.DB.ProductCollection.FindOne(ctx, bson.M{"name": product.ProductName}).Decode(&existingProduct)

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product already exists"})
		c.Abort()
		return
	}

	// Create the new product
	result, err := database.DB.ProductCollection.InsertOne(ctx, product)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully", "id": result.InsertedID})
}

/*
GetProductController retrieves a product from the database by its ID.

	It takes a product ID as input and returns a Product object and an error.
	If the product is not found in the database, it returns a "not found" error.
*/
func GetProductController(c *gin.Context) {
	productID := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(productID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var product productModel.Product
	err = database.DB.ProductCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func UpdateProductController(c *gin.Context) {}

/*
DeleteProductController deletes a specific product by ID.

	It takes a product ID as input and returns an error if the product is not found in the database.
	It deletes the product and returns a success message.
*/
func DeleteProductController(c *gin.Context) {
	productID := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// check if the product exists first
	errFind := database.DB.ProductCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&productModel.Product{})

	if errFind != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product is not exist"})
		c.Abort()
		return
	}
	_, err = database.DB.ProductCollection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		c.Abort()

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
