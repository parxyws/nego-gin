package route

import (
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/product"
)

func ProductRoute(route *gin.RouterGroup, controller product.ProductController) {
	products := route.Group("/products")
	{
		products.GET("/", controller.ListProducts)
		products.GET("/featured", controller.ListSellerProducts) // Assuming this for now if no specific featured handler
		products.GET("/:productId", controller.GetProductDetails)
		products.GET("/slug/:slug", controller.GetProductSlug)
	}

	seller := route.Group("/sellers/products")
	{
		seller.POST("/", controller.CreateProduct)
		seller.PUT("/:productId", controller.UpdateProduct)
		seller.DELETE("/:productId", controller.DeleteProduct)
		seller.PATCH("/:productId/status", controller.UpdateProductStatus)
		seller.GET("/", controller.ListSellerProducts)
	}
}
