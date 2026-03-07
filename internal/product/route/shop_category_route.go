package route

import (
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/product"
)

func ShopCategoryRoute(route *gin.RouterGroup, controller product.ShopCategoryController, authMiddleware *jwt.GinJWTMiddleware) {
	// Seller-protected shop category management
	sellerCategories := route.Group("/sellers/shop-categories")
	sellerCategories.Use(authMiddleware.MiddlewareFunc())
	{
		sellerCategories.GET("", controller.ListShopCategories)
		sellerCategories.GET("/:shopCategoryId", controller.GetShopCategory)
		sellerCategories.POST("", controller.CreateShopCategory)
		sellerCategories.PUT("/:shopCategoryId", controller.UpdateShopCategory)
		sellerCategories.DELETE("/:shopCategoryId", controller.DeleteShopCategory)
	}

	// Public shop browsing
	shopCategories := route.Group("/shops/:sellerId/categories")
	{
		shopCategories.GET("", controller.ListShopCategoriesPublic)
		shopCategories.GET("/:shopCategoryId", controller.GetShopCategoryPublic)
	}
}
