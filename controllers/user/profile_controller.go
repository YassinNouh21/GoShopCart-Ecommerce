package user

import (
	"context"
	"errors"
	"net/http"
	"time"

	"ecommerce/database"
	userModels "ecommerce/models/user"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrUnauthorized    = errors.New("Unauthorized")
	ErrInvalidID       = errors.New("Invalid id")
	ErrUserNotFound    = errors.New("User not found")
	ErrInvalidRequest  = errors.New("Invalid request body")
	ErrUpdateFailed    = errors.New("Failed to update user")
	ErrProductNotFound = errors.New("Product not found")
)

// profileRespose represents the response structure for the profile request.
type profileRespose struct {
	UserID         string               `json:"userid"`
	FirstName      string               `json:"first_name"`
	LastName       string               `json:"last_name"`
	Email          string               `json:"email"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
	AddressDetails []userModels.Address `json:"address"`
}

/*
GetProfileController handles the retrieval of user profile information.

	It retrieves the user ID from the request context, queries the database for the user with the corresponding ID,
	and returns the user's profile information.

Possible Errors:
  - ErrUnauthorized: If the user ID is not found in the request context.
  - ErrInvalidID: If the user ID in the request context is not a valid ObjectID.
  - ErrUserNotFound: If no user with the provided ID exists in the database.
*/
func GetProfileController(c *gin.Context) {
	// Create a context with a timeout of 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Retrieve the user ID from the request context
	userID, errBool := c.Get("user_id")
	if !errBool {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ErrUnauthorized.Error()})
		c.Abort()
		return
	}

	// Convert the user ID to an ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidID.Error()})
		c.Abort()
		return
	}

	// Query the database for the user with the corresponding ID
	filter := bson.M{"_id": objectID}
	var user userModels.User
	err = database.DB.UserCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrUserNotFound.Error()})
		c.Abort()
		return
	}

	// Create a profile response object and return it as JSON
	userUpdated := profileRespose{
		UserID:         user.ID.Hex(),
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Email:          user.Email,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
		AddressDetails: user.AddressDetails,
	}
	c.IndentedJSON(http.StatusOK, userUpdated)
}

// UpdateProfile represents the request structure for updating user profile information.
type UpdateProfile struct {
	FirstName      string               `json:"first_name" validate:"required"`
	LastName       string               `json:"last_name" validate:"required"`
	Email          string               `json:"email" validate:"required,email"`
	AddressDetails []userModels.Address `json:"address" validate:"required"`
}

/*
UpdateProfileController handles the updating of user profile information.

It retrieves the user ID from the request context, parses the JSON request body containing the updated user information,
validates the request body, updates the user's profile information in the database, and returns a success message.

Possible Errors:
  - ErrUnauthorized: If the user ID is not found in the request context.
  - ErrInvalidRequest: If the request body is not in the expected format or contains invalid data.
  - ErrUpdateFailed: If an error occurs while updating the user's profile information in the database.
  - ErrUserNotFound: If no user with the provided ID exists in the database.
*/
func UpdateProfileController(c *gin.Context) {
	// Create a context with a timeout of 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Retrieve the user ID from the request context
	userID, isFound := c.Get("user_id")
	if !isFound {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ErrUnauthorized.Error()})
		c.Abort()
		return
	}

	// Parse the JSON request body containing the updated user information
	var updatedUser UpdateProfile
	err := c.ShouldBindJSON(&updatedUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidRequest.Error()})
		c.Abort()
		return
	}

	// Convert the user ID to an ObjectID and create a filter for the update query
	objectID, err := primitive.ObjectIDFromHex(userID.(string))
	filter := bson.M{"_id": objectID}

	// Update the user's profile information in the database
	update := bson.M{"$set": updatedUser}
	result, err := database.DB.UserCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrUpdateFailed.Error()})
		c.Abort()
		return
	}

	// Check if the update query matched any documents in the database
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrUserNotFound.Error()})
		c.Abort()
		return
	}

	// Return a success message
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
