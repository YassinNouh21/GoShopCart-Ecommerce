package middlewares

/*

The Authentication middleware function, which validates the user's authentication token.

Functions:
- Authentication: Validates the user's authentication token.
*/

import (
	"ecommerce/helpers"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

// Authentication is a middleware function that validates the user's authentication token.
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		authHeader := c.GetHeader("Authorization")
		log.Println("midlleware: authHeader", authHeader)
		clientToken = strings.TrimPrefix(authHeader, "Bearer ")
		if clientToken == "" {
			c.JSON(401, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}
		log.Println("midlleware: clientToken", clientToken[0:5])

		userClaim, err := helpers.ValidateToken(clientToken)
		log.Println("midlleware: userclaim", userClaim)

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
