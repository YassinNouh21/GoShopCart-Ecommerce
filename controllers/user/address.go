package user

import (
	"context"
	"fmt"
	"github.com/YassinNouh21/GoShopCart-Ecommerce/database"
	"github.com/YassinNouh21/GoShopCart-Ecommerce/models/user"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
AddAddressController handles the creation of a new address for a user.

	It retrieves the user ID from the request context, retrieves the user from the database,
	validates the request body, generates a new address ID, updates the user with the new address, and returns a success message.

Possible Errors:
  - Unauthorized: If the user ID is not found in the request context.
  - User not found: If the user with the provided ID is not found in the database.
  - Invalid request body: If the request body is not in the expected format or contains invalid data.
  - Failed to create address: If an error occurs while updating the user with the new address.
*/
func AddAddressController(c *gin.Context) {

	userID, errBool := c.Get("user_id")
	if !errBool {
		c.JSON(401, gin.H{"error": "Unauthorized"})
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
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		c.Abort()
		return
	}

	// Bind the request body to the Address model
	var address user.Address

	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	validator := validator.New()

	err = validator.Struct(address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	addressID := primitive.NewObjectID()
	address.AddressID = addressID
	// create slides array
	update := bson.M{
		"$push": bson.M{
			"address_details": address,
		},
	} // Update the user in the database
	_, err = database.DB.UserCollection.UpdateOne(ctx, primitive.M{"_id": userID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create address"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Address created successfully"})
}

/*
DeleteAllAddressController handles the deletion of all addresses for a user.
	It retrieves the user ID from the request context, retrieves the user from the database,
	removes all addresses from the user's address details field, updates the user in the database, and returns a success message.
Possible Errors:
	- Unauthorized: If the user ID is not found in the request context.
	- User not found: If the user with the provided ID is not found in the database.
	- Error deleting address: If an error occurs while deleting the addresses.
*/

func DeleteAllAddressController(c *gin.Context) {
	userID, errBool := c.Get("user_id")
	if !errBool {
		c.JSON(401, gin.H{"error": "Unauthorized"})
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
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		c.Abort()
		return
	}
	update := bson.M{"$set": bson.M{"address_details": []user.Address{}}}
	filter := bson.M{"_id": userID}

	_, err = database.DB.UserCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Error deleting address"})
		c.Abort()
		return
	}
	err = database.DB.UserCollection.FindOne(ctx, primitive.M{"_id": userID}).Decode(&existingUser)

	// create slides array

	c.JSON(http.StatusOK, gin.H{"message": "All Addresses successfully deleted"})
}

/*
DeleteAddressWithIdController handles the deletion of a specific address for a user.
	It retrieves the user ID from the request context, retrieves the address ID from the request parameters,
	retrieves the user from the database, removes the matching address from the user's address details field,
	updates the user in the database, and returns a success message.
Possible Errors:
	- Unauthorized: If the user ID is not found in the request context.
	- There is no address id provided: If the address ID is not provided in the request parameters.
	- User not found: If the user with the provided ID is not found in the database.
	- Invalid address ID: If the provided address ID is not a valid MongoDB ObjectID.
	- Address not found: If the user's address details field does not contain an address with the provided ID.
*/

func DeleteAddressWithIdController(c *gin.Context) {

	userID, errBool := c.Get("user_id")

	addressId := c.Param("address_id")
	if !errBool {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	if addressId == "" {
		c.JSON(401, gin.H{"error": "There is no address id provided"})
		c.Abort()
		return
	}
	// Get the user from the database
	userID, _ = primitive.ObjectIDFromHex(userID.(string))

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
	defer cancel()
	addressIDObj, err := primitive.ObjectIDFromHex(addressId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	// Define the filter to match the user and the address ID
	filter := bson.M{
		"_id":                 userID,
		"address_details._id": addressIDObj,
	}

	// Define the update to remove the matching address
	update := bson.M{
		"$pull": bson.M{
			"address_details": bson.M{"_id": addressIDObj},
		},
	}

	updateResult, err := database.DB.UserCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		c.Abort()
		return
	}
	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		c.Abort()
		return
	}
	if updateResult.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		c.Abort()
		return
	}
	// create slides array
	message := fmt.Sprintf("Address with ID %s deleted successfully", addressId)

	c.JSON(http.StatusOK, gin.H{"message": message})
}

/*
GetAddressController retrieves all addresses for a user.

	It retrieves the user ID from the request context, retrieves the user from the database,
	and returns the user's address details field.

Possible Errors:
  - Unauthorized: If the user ID is not found in the request context.
  - User not found: If the user with the provided ID is not found in the database.
*/
func GetAddressController(c *gin.Context) {

	userID, errBool := c.Get("user_id")
	if !errBool {
		c.JSON(401, gin.H{"error": "Unauthorized"})
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
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		c.Abort()
		return
	}
	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		c.Abort()
		return
	}
	// create slides array

	c.IndentedJSON(http.StatusOK, gin.H{"message": existingUser.AddressDetails})
}
