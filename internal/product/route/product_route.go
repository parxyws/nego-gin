package route

import (
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/product"
)

func ProductRoute(route *gin.RouterGroup, controller product.ProductController, authMiddleware *jwt.GinJWTMiddleware) {
	// Public product browsing
	products := route.Group("/products")
	{
		products.GET("", controller.ListProducts)
		products.GET("/featured", controller.ListFeaturedProducts)
		products.GET("/:productId", controller.GetProductDetails)
		products.GET("/slug/:slug", controller.GetProductSlug)
	}

	// Seller-protected product management
	sellerProducts := route.Group("/sellers/products")
	sellerProducts.Use(authMiddleware.MiddlewareFunc())
	{
		sellerProducts.GET("", controller.ListSellerProducts)
		sellerProducts.POST("", controller.CreateProduct)
		sellerProducts.PUT("/:productId", controller.UpdateProduct)
		sellerProducts.DELETE("/:productId", controller.DeleteProduct)
		sellerProducts.PATCH("/:productId/status", controller.UpdateProductStatus)
	}
}
