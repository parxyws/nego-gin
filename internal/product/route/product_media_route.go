package route

import (
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/product"
)

func ProductMediaRoute(route *gin.RouterGroup, controller product.ProductMediaController, authMiddleware *jwt.GinJWTMiddleware) {
	// Public — list media
	products := route.Group("/products")
	{
		products.GET("/:productId/media", controller.ListProductMedia)
	}

	// Seller-protected — manage media
	sellerProducts := route.Group("/sellers/products")
	sellerProducts.Use(authMiddleware.MiddlewareFunc())
	{
		sellerProducts.POST("/:productId/media", controller.UploadMedia)
		sellerProducts.PUT("/:productId/media/:mediaId", controller.UpdateProductMedia)
		sellerProducts.DELETE("/:productId/media/:mediaId", controller.DeleteProductMedia)
	}
}
