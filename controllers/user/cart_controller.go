package user

import (
	"context"
	"errors"
	"fmt"
	database "github.com/YassinNouh21/GoShopCart-Ecommerce/database"
	"github.com/YassinNouh21/GoShopCart-Ecommerce/models/user"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	// ErrFailedUpdate is returned when the cart update fails.
	ErrFailedUpdate = errors.New("Failed to update cart")

	// ErrInvalidCardId is returned when an invalid cart ID is provided.
	ErrInvalidCardId = errors.New("Invalid cart ID")

	// ErrCartNotFound is returned when the cart is not found.
	ErrCartNotFound = errors.New("This cart is not found")

	// ErrCartNotUpdated is returned when the cart cannot be updated.
	ErrCartNotUpdated = errors.New("Failed to update cart")

	// ErrCartNotCreate is returned when the cart cannot be created.
	ErrCartNotCreate = errors.New("Failed to create cart")

	// ErrCartIdNotProvided is returned when the cart ID is not provided in the request body.
	ErrCartIdNotProvided = errors.New("Cannot provide cartID in the request body")
)

/*
	GetCartController returns a cart for the authenticated user.

	Possible errors:
		- ErrUnauthorized: if the user is not authenticated
		- ErrUserNotFound: if the user cannot be found in the database
		- ErrInvalidCartID: if the cart ID in the request body is invalid
		- ErrCartNotFound: if the cart cannot be found in the database
*/

func GetCartController(c *gin.Context) {

	userID, errBool := c.Get("user_id")
	if !errBool {
		c.JSON(401, gin.H{"error": ErrUnauthorized.Error()})
		c.Abort()
		return
	}

	// Get the user from the database
	var existingUser user.User
	userID, _ = primitive.ObjectIDFromHex(userID.(string))
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
	err := database.DB.UserCollection.FindOne(ctx, primitive.M{"_id": userID}).Decode(&existingUser)
	defer cancel()

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrUserNotFound.Error()})
		c.Abort()
		return
	}
	log.Println("posterror", err, "user", existingUser)
	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": ErrUserNotFound.Error()})
		c.Abort()
		return
	}
	// create slides array

	c.IndentedJSON(http.StatusOK, gin.H{"message": existingUser.UserCart})
}

func AddCartController(c *gin.Context) {

	userID, errBool := c.Get("user_id")
	if !errBool {
		c.JSON(401, gin.H{"error": ErrUnauthorized.Error()})
		c.Abort()
		return
	}

	// Get the user from the database
	var existingUser user.User
	userID, _ = primitive.ObjectIDFromHex(userID.(string))
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
	defer cancel()
	err := database.DB.UserCollection.FindOne(ctx, primitive.M{"_id": userID}).Decode(&existingUser)
	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": ErrUserNotFound.Error()})
		c.Abort()
		return
	}

	var cart user.Cart

	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	validator := validator.New()
	err = validator.Struct(&cart)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	count, err := database.DB.ProductCollection.CountDocuments(ctx, bson.M{"_id": cart.ProductID})
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrProductNotFound.Error()})
		c.Abort()
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	filterSearchProductIdCart := bson.M{
		"_id":                  userID,
		"user_cart.product_id": cart.ProductID,
	}
	// check if the product is already in the cart
	count, err = database.DB.UserCollection.CountDocuments(ctx, filterSearchProductIdCart)
	log.Println("number in cart", err, "cart", count)
	if count > 0 {
		// increase the quantity and update the cart
		filterSearchProductIdCart := bson.M{
			"_id":                  userID,
			"user_cart.product_id": cart.ProductID,
		}
		updateIncrement := bson.M{
			"$inc": bson.M{
				"user_cart.$.quantity": cart.Quantity,
			},
			"$set": bson.M{
				"user_cart.$.updated_at": time.Now().UTC(),
			},
		}
		updated, err := database.DB.UserCollection.UpdateOne(ctx, filterSearchProductIdCart, updateIncrement)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if updated.MatchedCount == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": ErrFailedUpdate.Error()})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Cart updated successfully"})
	} else {

		log.Println("error", err, "cart", cart)
		cartID := primitive.NewObjectID()
		cart.CartID = cartID
		cart.CreatedAt = time.Now().UTC()
		cart.UpdatedAt = time.Now().UTC()

		// create slides array
		update := bson.M{
			"$push": bson.M{
				"user_cart": &cart,
			},
		} // Update the user in the database
		_, err = database.DB.UserCollection.UpdateOne(ctx, primitive.M{"_id": &userID}, update)
		if err != nil {
			log.Println("posterror", err, "address", cart)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
			return
		}
		message := fmt.Sprintf("Cart with ID %s created successfully", cart.CartID.Hex())
		c.JSON(http.StatusOK, gin.H{"message": message})
	}
}

func DeleteAllCartController(c *gin.Context) {

	userID, errBool := c.Get("user_id")
	if !errBool {
		c.JSON(401, gin.H{"error": ErrUnauthorized.Error()})
		c.Abort()
		return
	}

	// Get the user from the database
	var existingUser user.User
	userID, _ = primitive.ObjectIDFromHex(userID.(string))
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
	defer cancel()
	err := database.DB.UserCollection.FindOne(ctx, primitive.M{"_id": userID}).Decode(&existingUser)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrUserNotFound.Error()})
		c.Abort()
		return
	}
	log.Println("posterror", err, "user", existingUser)
	update := bson.M{"$set": bson.M{"user_cart": []user.Cart{}}}
	filter := bson.M{"_id": userID}

	_, err = database.DB.UserCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Error deleting carts"})
		c.Abort()
		return
	}

	// create slides array

	c.JSON(http.StatusOK, gin.H{"message": "All carts are successfully deleted"})
}

// func DeleteAddressWithIdController(c *gin.Context) {

// 	userID, errBool := c.Get("user_id")

// 	addressId := c.Param("address_id")
// 	log.Println("addressId", addressId, c.Params)
// 	if !errBool {
// 		c.JSON(401, gin.H{"error": "Unauthorized"})
// 		c.Abort()
// 		return
// 	}
// 	if addressId == "" {
// 		c.JSON(401, gin.H{"error": "There is no address id provided"})
// 		c.Abort()
// 		return
// 	}
// 	// Get the user from the database
// 	var existingUser user.User
// 	userID, _ = primitive.ObjectIDFromHex(userID.(string))

// 	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
// 	defer cancel()
// 	addressIDObj, err := primitive.ObjectIDFromHex(addressId)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
// 		return
// 	}

// 	// Define the filter to match the user and the address ID
// 	filter := bson.M{
// 		"_id":                userID,
// 		"addressdetails._id": addressIDObj,
// 	}

// 	// Define the update to remove the matching address
// 	update := bson.M{
// 		"$pull": bson.M{
// 			"addressdetails": bson.M{"_id": addressIDObj},
// 		},
// 	}

// 	updateResult, err := database.DB.UserCollection.UpdateOne(ctx, filter, update)

// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 		c.Abort()
// 		return
// 	}
// 	log.Println("posterror", err, "user", existingUser)
// 	if err != nil {

// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 		c.Abort()
// 		return
// 	}
// 	if updateResult.MatchedCount == 0 {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
// 		c.Abort()
// 		return
// 	}
// 	// create slides array
// 	message := fmt.Sprintf("Address with ID %s deleted successfully", addressId)

//		c.JSON(http.StatusOK, gin.H{"message": message})
//	}
func UpdateCartController(c *gin.Context) {
	userID, errBool := c.Get("user_id")
	cartId := c.Param("cart_id")
	log.Println("cartId", cartId, c.Params)
	cartIdObj, err := primitive.ObjectIDFromHex(cartId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidCardId.Error()})
		return
	}
	if !errBool {
		c.JSON(401, gin.H{"error": ErrUnauthorized.Error()})
		c.Abort()
		return
	}

	// Get the user from the database
	var existingUser user.User
	userID, _ = primitive.ObjectIDFromHex(userID.(string))
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
	defer cancel()
	err = database.DB.UserCollection.FindOne(ctx, primitive.M{"_id": userID}).Decode(&existingUser)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrUserNotFound.Error()})
		c.Abort()
		return
	}

	var cart user.CartWithoutId

	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	if cart.CartID != primitive.NilObjectID {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrCartIdNotProvided.Error()})
		return
	}
	validator := validator.New()
	err = validator.Struct(&cart)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	// check if cart exist

	filterIdCart := bson.M{
		"_id":           &userID,
		"user_cart._id": &cartIdObj,
	}
	log.Println("filterIdCart", filterIdCart, "updateCartWithID", filterIdCart)

	count, err := database.DB.UserCollection.CountDocuments(ctx, &filterIdCart)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrCartNotFound.Error()})
		c.Abort()
		return
	}
	// chekc if the product is already in the cart

	log.Println("error", err, "cart", cart)
	cart.UpdatedAt = time.Now().UTC()

	// create slides array except the _id
	cart.CartID = cartIdObj
	updateCartWithID := bson.M{
		"$set": bson.M{
			"user_cart.$": &cart,
		},
	}
	updated, err := database.DB.UserCollection.UpdateOne(ctx, filterIdCart, updateCartWithID)
	if updated.MatchedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrCartNotUpdated.Error()})
		c.Abort()
		return
	}
	if err != nil {
		log.Println("posterror", err, "address", cart)
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrCartNotCreate.Error()})
		return
	}
	message := fmt.Sprintf("Cart with ID %s updated successfully", cart.CartID.Hex())
	c.JSON(http.StatusOK, gin.H{"message": message})
}
