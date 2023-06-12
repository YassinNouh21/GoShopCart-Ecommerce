package routes

import (
	"github.com/YassinNouh21/GoShopCart-Ecommerce/controllers/auth"

	"github.com/gin-gonic/gin"
)

/*
	Functions:
	- GetAuthRoutes: Sets up the authentication routes for user authentication.
*/

// GetAuthRoutes sets up the authentication routes for user authentication.
func GetAuthRoutes(userRoutes *gin.RouterGroup) {
	userRoutes.POST("/signin", auth.SignInController)
	userRoutes.POST("/signup", auth.SignUpController)
	userRoutes.POST("/tokenrefresh", auth.TokenRefreshController)
}
