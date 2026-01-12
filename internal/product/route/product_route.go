package route

import (
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/product"
)

func ProductRoute(route *gin.RouterGroup, controller product.ProductController) {
	products := route.Group("/products")
	{
		products.GET("/")
		products.GET("/featured")
		products.GET("/:productId")
		products.GET("/slug/:slug")
	}

	seller := route.Group("/sellers/products")
	{
		seller.POST("/")
	}
}
