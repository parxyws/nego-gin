package route

import (
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/product"
)

func CategoryRoute(route *gin.RouterGroup, controller product.CategoryController, authMiddleware *jwt.GinJWTMiddleware) {
	// Public category browsing
	categories := route.Group("/categories")
	{
		categories.GET("", controller.ListCategories)
		categories.GET("/:categoryId", controller.GetCategory)
		categories.GET("/:categoryId/products", controller.ListProductsInCategory)
	}

	// Admin-only category management
	adminCategories := route.Group("/admin/categories")
	adminCategories.Use(authMiddleware.MiddlewareFunc())
	{
		adminCategories.POST("", controller.CreateCategory)
		adminCategories.PUT("/:categoryId", controller.UpdateCategory)
		adminCategories.DELETE("/:categoryId", controller.RemoveCategory)
	}
}
