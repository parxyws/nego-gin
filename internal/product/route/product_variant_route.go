package route

import (
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/product"
)

func ProductVariantRoute(route *gin.RouterGroup, controller product.ProductVariantsController, authMiddleware *jwt.GinJWTMiddleware) {
	// Public — list variants
	products := route.Group("/products")
	{
		products.GET("/:productId/variants", controller.ListProductVariants)
	}

	// Seller-protected — manage variants
	sellerProducts := route.Group("/sellers/products")
	sellerProducts.Use(authMiddleware.MiddlewareFunc())
	{
		sellerProducts.POST("/:productId/variants", controller.CreateProductVariant)
		sellerProducts.PUT("/:productId/variants/:variantId", controller.UpdateProductVariant)
		sellerProducts.DELETE("/:productId/variants/:variantId", controller.DeleteProductVariant)
	}
}
