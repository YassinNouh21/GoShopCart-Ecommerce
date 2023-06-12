package auth

import (
	goContext "context"
	"errors"
	"fmt"
	"github.com/YassinNouh21/GoShopCart-Ecommerce/database"
	helpers "github.com/YassinNouh21/GoShopCart-Ecommerce/helpers"
	userModel "github.com/YassinNouh21/GoShopCart-Ecommerce/models/user"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	errInvalidRequestBody = errors.New("invalid request body")
	errUserAlreadyExists  = errors.New("user with that email already exists")
	errUserNotFound       = errors.New("user not found")
	errHashingPassword    = errors.New("error while hashing password")
	errIncorrectPassword  = errors.New("password is incorrect")
	errGeneratingToken    = errors.New("error while generating token")
	errUpdatingToken      = errors.New("error while updating token")
	errUserNotFoundByID   = errors.New("user not found with this ID")
)

/*
SignUpController handles the user registration process.

It parses the JSON request body into a user model, validates the request body, checks if the user already exists, generates a token, and creates a new user record in the database.

Errors:
	- Invalid request body: If the request body is not in the expected format or contains invalid data.
	- User already exists: If a user with the provided email already exists in the database.
	- Error while generating token: If an error occurs while generating the authentication token.
	- Error while inserting user: If an error occurs while inserting the new user record into the database.
*/

func SignUpController(context *gin.Context) {
	ctx, cancel := goContext.WithTimeout(goContext.Background(), 10*time.Second)
	defer cancel()

	var user userModel.User
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": errInvalidRequestBody.Error(),
		})
		return
	}

	if err := validator.New().Struct(user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": errInvalidRequestBody.Error(),
		})
		return
	}

	count, err := database.DB.UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if count > 0 {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": errUserAlreadyExists.Error(),
		})
		return
	}

	timeAt, err := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	timeUpdateted, err := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.CreatedAt = timeAt
	user.UpdatedAt = timeUpdateted
	user.Password, _ = helpers.HashPassword(user.Password)

	user.ID = primitive.NewObjectID()
	userID := user.ID.Hex()
	userClaims := helpers.CreateUserClaims(user.Email, user.FirstName, userID)
	tokenGenerated, refreshToken, err := helpers.GenerateToken(*userClaims)
	user.Token = tokenGenerated
	user.RefreshToken = refreshToken
	if err != nil {
		err := fmt.Sprintf("Error while generating token: %s", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	user.AddressDetails = []userModel.Address{}
	user.OrderStatus = []userModel.Order{}
	user.UserCart = []userModel.Cart{}

	_, err = database.DB.UserCollection.InsertOne(ctx, user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	} else {
		context.JSON(http.StatusOK, gin.H{
			"message": "User created successfully",
		})
		return
	}
}

// SignInResponse represents the response structure for the signing request.
type SignInResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

/*
SignInController handles the user login process.
	It parses the JSON request body into a user model, validates the request body, retrieves the user from the database, verifies the password, generates a new access token, and updates the user's tokens.

Errors:
	- Invalid request body: If the request body is not in the expected format or contains invalid data.
	- User not found: If the user with the provided email does not exist in the database.
	- Error while hashing password: If an error occurs while comparing the passwords.
	- Error while generating token: If an error occurs while generating the authentication token.
	- Error while updating token: If an error occurs while updating the user's access and refresh tokens.
	- Password is incorrect: If the provided password does not match the user's stored password.
*/

func SignInController(context *gin.Context) {
	ctx, cancel := goContext.WithTimeout(goContext.Background(), 30*time.Second)
	defer cancel()

	var user userModel.User
	var loginUser userModel.User

	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": errInvalidRequestBody.Error(),
		})
		return
	}
	if err := validator.New().Struct(user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": errInvalidRequestBody.Error(),
		})
		return
	}

	err := database.DB.UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&loginUser)
	defer cancel()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": errUserNotFound.Error(),
		})
		return
	}

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": errHashingPassword.Error(),
		})
		return
	}

	isValid := helpers.VerifyPassword(loginUser.Password, user.Password)
	userId := loginUser.ID.Hex()
	userClaim := *helpers.CreateUserClaims(loginUser.Email, loginUser.FirstName, userId)
	accessToken, refreshToken, err := helpers.GenerateToken(userClaim)
	defer cancel()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": errGeneratingToken.Error(),
		})
		return
	}

	err = helpers.UpdateToken(accessToken, refreshToken, loginUser.ID)
	if !isValid {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": errIncorrectPassword.Error(),
		})
		return
	}

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": errUpdatingToken.Error(),
		})
		return
	}

	signInRes := SignInResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	context.JSON(http.StatusOK, gin.H{
		"message": signInRes,
	})
}

/*
GetUserIdController handles the retrieval of user information by user ID.

	It retrieves the user ID from the request parameters, queries the database for the user with the corresponding ID, and returns a response indicating whether the user was found.

Errors:
  - User not found: If no user with the provided ID exists in the database.
*/
func GetUserIdController(context *gin.Context) {
	userId := context.Param("user_id")
	ctx, cancel := goContext.WithTimeout(goContext.Background(), 10*time.Second)
	defer cancel()

	err := database.GetCollectionMongoDB("users").FindOne(ctx, bson.M{"_id": userId}).Decode(&userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("%s: %s", errUserNotFoundByID.Error(), userId),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "User found",
	})
}

// TokenRefreshResponse represents the request body for refreshing the access token.
type TokenRefreshResponse struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

/*
TokenRefreshController handles the token refresh process.

	It parses the JSON request body containing the refresh token, validates the request body, generates a new access token using the refresh token, and returns the new access token.

Errors:
  - Invalid request body: If the request body is not in the expected format or contains invalid data.
  - Error while generating token: If an error occurs while generating the new authentication token.
*/
func TokenRefreshController(context *gin.Context) {
	var refreshToken TokenRefreshResponse

	if err := context.ShouldBindJSON(&refreshToken); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": errInvalidRequestBody.Error(),
		})
		return
	}

	if err := validator.New().Struct(refreshToken); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": errInvalidRequestBody.Error(),
		})
		return
	}
	// request a new access token
	accessToken, err := helpers.GenerateNewAccessToken(refreshToken.RefreshToken)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": accessToken,
	})
}
