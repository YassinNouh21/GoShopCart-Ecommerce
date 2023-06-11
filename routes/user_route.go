package routes

import (
	"ecommerce/controllers/user"
	"github.com/gin-gonic/gin"
)

// AddressRoutes sets up the address routes of the user.
func AddressRoutes(addressRoutes *gin.RouterGroup) {
	addressRoutes.GET("/address", user.GetAddressController)
	addressRoutes.POST("/address", user.AddAddressController)
	addressRoutes.DELETE("/address", user.DeleteAllAddressController)
	addressRoutes.DELETE("/address/:address_id", user.DeleteAddressWithIdController)
}
