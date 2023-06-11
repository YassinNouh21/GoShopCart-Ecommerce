package helpers

import (
	"context"
	"ecommerce/database"
	userModel "ecommerce/models/user"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
	This package implements functions for generating, validating, and updating JWT tokens for user authentication.
	It also includes functions for generating new access tokens based on refresh tokens and performing token validation.

	Error Handling:
	This package defines the following errors:
	- "User token is not updated": Returned when the token update operation fails.
	- "Token is expired": Returned when a token is expired and cannot be validated.
	- "Invalid user ID": Returned when the provided user ID is not a valid MongoDB ObjectID.
	- "User not found": Returned when a user is not found in the database during token validation.
	- "Token is not valid": Returned when a token is not valid during validation or does not match the user's token in the database.
	- "Error while parsing claims": Returned when there is an error while parsing JWT claims.
	- "Error while generating new token": Returned when there is an error during the generation of a new token.
	- "Error while updating token": Returned when there is an error updating the user's token in the database.

*/

// mongoDBCollectionUser represents the MongoDB collection for user data.
var mongoDBCollectionUser *mongo.Collection = database.DB.UserCollection

// UserClaims represents the custom claims for a JWT token.
type UserClaims struct {
	Email     string
	FirstName string
	ID        string
	jwt.StandardClaims
}

// CreateUserClaims creates a new UserClaims instance with the provided email, first name, and ID.
func CreateUserClaims(email string, firstName string, id string) *UserClaims {
	return &UserClaims{
		Email:     email,
		FirstName: firstName,
		ID:        id,
	}
}

// API_KEY is the secret key used for JWT token generation and verification.
var API_KEY = []byte(os.Getenv("SECRET_JWT"))

// GenerateToken generates a new JWT token and refresh token based on the provided user claims.
// It returns the signed token, signed refresh token, and any error encountered.
func GenerateToken(userclaim UserClaims) (signedToken string, signedRefreshToken string, err error) {
	secretKey := API_KEY

	// Set expiration time for the token
	userclaim.StandardClaims = jwt.StandardClaims{
		ExpiresAt: time.Now().Local().Add(time.Minute * time.Duration(5)).Unix(),
	}

	// Create a new token with the user claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userclaim)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Fatalln("1", err)
		return "", "", err
	}

	// Set expiration time for the refresh token
	userclaim.StandardClaims = jwt.StandardClaims{
		ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(2190)).Unix(),
	}

	// Create a new refresh token with the user claims
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, userclaim)

	// Sign the refresh token with the secret key
	refreshTokenString, err := refreshToken.SignedString(secretKey)
	if err != nil {
		log.Fatalln("2", err)
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}

// UpdateToken updates the user's access token and refresh token in the database.
// It takes the signed access token, signed refresh token, and user ID as parameters.
// It returns an error if the token update operation fails.
func UpdateToken(signedToken string, signedRefreshToken string, userId primitive.ObjectID) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

	// Create the update operation with the new tokens and update time
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "token", Value: signedToken}, {Key: "refreshtoken", Value: signedRefreshToken}, {Key: "updateat", Value: time.Now()}}}}

	// Create the filter for the user ID
	filter := bson.D{{Key: "_id", Value: userId}}

	// Perform the update operation on the user collection
	updatedUser, err := database.DB.UserCollection.UpdateOne(ctx, filter, update)
	defer cancel()

	if updatedUser.MatchedCount == 0 {
		return errors.New("User token is not updated")
	}

	if err != nil {
		log.Panic(err)
		return err
	}

	return nil
}

// ValidateToken validates the provided JWT token and returns the claims if valid.
// It also checks the token against the user's token in the database for additional validation.
// It returns the claims and an error message if any issue occurs during validation.
func ValidateToken(verifyToken string) (claim *UserClaims, errorMessage string) {
	token, err := jwt.ParseWithClaims(verifyToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(API_KEY), nil
	})
	fmt.Println("token", err, verifyToken)

	// Check if the token is expired
	if err != nil && strings.Contains(err.Error(), "expired") {
		return nil, "token is expired"
	}

	claim, ok := token.Claims.(*UserClaims)
	userIdPrimitive, err := primitive.ObjectIDFromHex(claim.ID)
	if err != nil {
		return nil, "invalid user id"
	}

	// Check if the user exists in the database
	var user userModel.User
	err = mongoDBCollectionUser.FindOne(context.Background(), bson.M{"_id": &userIdPrimitive}).Decode(&user)
	if err != nil {
		return nil, "user not found"
	}

	// Check if the token matches the user's token in the database
	if user.Token != verifyToken {
		return nil, "token is not valid"
	}

	if !ok {
		return nil, "error while parsing claims"
	}

	if !token.Valid {
		return nil, "token is not valid"
	}

	return claim, ""
}

// ValidateRefreshToken validates the provided refresh token and returns the claims if valid.
// It also checks the token against the user's refresh token in the database for additional validation.
// It returns the claims and an error message if any issue occurs during validation.
func ValidateRefreshToken(verifyToken string) (claim *UserClaims, errorMessage string) {
	token, err := jwt.ParseWithClaims(verifyToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(API_KEY), nil
	})
	fmt.Println("refresh token", err, verifyToken)

	// Check if the token is expired
	if err != nil && strings.Contains(err.Error(), "expired") {
		return nil, "token is expired"
	}

	claim, ok := token.Claims.(*UserClaims)
	userIdPrimitive, err := primitive.ObjectIDFromHex(claim.ID)
	if err != nil {
		return nil, "invalid user id"
	}

	// Check if the user exists in the database
	var user userModel.User
	err = mongoDBCollectionUser.FindOne(context.Background(), bson.M{"_id": &userIdPrimitive}).Decode(&user)
	if err != nil {
		return nil, "user not found"
	}

	// Check if the token matches the user's refresh token in the database
	if user.RefreshToken != verifyToken {
		return nil, "token is not valid"
	}

	if !ok {
		return nil, "error while parsing claims"
	}

	if !token.Valid {
		return nil, "token is not valid"
	}

	return claim, ""
}

// GenerateNewAccessToken generates a new access token based on the provided refresh token.
// It validates the refresh token and checks the user's existence in the database.
// It returns the signed access token and any error encountered.
func GenerateNewAccessToken(refreshToken string) (signedToken string, err error) {
	claim, errString := ValidateRefreshToken(refreshToken)
	if errString != "" {
		return "", errors.New(errString)
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user userModel.User
	userIdPrimitive, err := primitive.ObjectIDFromHex(claim.ID)
	if err != nil {
		return "", errors.New("invalid user id")
	}

	err = mongoDBCollectionUser.FindOne(ctx, bson.M{"_id": &userIdPrimitive}).Decode(&user)
	if err != nil {
		return "", errors.New("user not found")
	}

	signedToken, _, err = GenerateToken(*claim)
	if err != nil {
		return "", errors.New("error while generating new token")
	}

	_, err = mongoDBCollectionUser.UpdateOne(ctx, bson.M{"_id": &userIdPrimitive}, bson.M{"$set": bson.M{"token": signedToken}})
	if err != nil {
		return "", errors.New("error while updating token")
	}

	return signedToken, nil
}
