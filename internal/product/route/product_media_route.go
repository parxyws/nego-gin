package route

import (
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/product"
)

func ProductMediaRoute(route *gin.RouterGroup, controller product.ProductMediaController) {
	products := route.Group("/products")
	{
		products.GET("/:productId/media", controller.ListProductMedia)
	}

	seller := route.Group("/sellers/products")
	{
		seller.POST("/:productId/media", controller.UploadMedia)
		seller.PUT("/:productId/media/:mediaId", controller.UpdateProductMedia)
		seller.DELETE("/:productId/media/:mediaId", controller.DeleteProductMedia)
	}
}
