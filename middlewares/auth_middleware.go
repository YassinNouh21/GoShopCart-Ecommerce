package middlewares

/*

The Authentication middleware function, which validates the user's authentication token.

Functions:
- Authentication: Validates the user's authentication token.
*/

import (
	"github.com/YassinNouh21/GoShopCart-Ecommerce/helpers"
	"strings"

	"github.com/gin-gonic/gin"
)

// Authentication is a middleware function that validates the user's authentication token.
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		authHeader := c.GetHeader("Authorization")
		clientToken = strings.TrimPrefix(authHeader, "Bearer ")
		if clientToken == "" {
			c.JSON(401, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		userClaim, err := helpers.ValidateToken(clientToken)

		if err != "" {
			c.JSON(401, gin.H{
				"error": err,
			})
			c.Abort()
			return
		}
		c.Set("user_id", userClaim.ID)
		c.Next()
	}
}
