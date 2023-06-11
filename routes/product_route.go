package routes

import (
	productController "ecommerce/controllers/product"
	"github.com/gin-gonic/gin"
)

// ProductRoutes sets up the routes for the product endpoints.
func ProductRoutes(productRoutes *gin.RouterGroup) {
	productRoutes.POST("/", productController.CreateProductController)
	productRoutes.GET("/:id", productController.GetProductController)
	productRoutes.PUT("/:id", productController.UpdateProductController)
	productRoutes.DELETE("/:id", productController.DeleteProductController)
}

// ProductFilterRoutes sets up the routes for the product filter endpoints.
func ProductFilterRoutes(productRoutes *gin.RouterGroup) {
	productRoutes.GET("/price", productController.GetProductsByPriceRangeController)
	productRoutes.GET("/price/:price", productController.GetProductsByPriceController)
	productRoutes.GET("/keyword", productController.GetProductsByKeyword)
}
