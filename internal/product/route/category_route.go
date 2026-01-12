package route

import (
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/product"
)

func CategoryRoute(route *gin.RouterGroup, controller product.CategoryController) {
	category := route.Group("/categories")
	{
		category.GET("/")
		category.GET("/:categoryId")
		category.GET("/:categoryId/products")
	}

	admin := route.Group("/admin/categories")
	{
		admin.POST("/")
		admin.PUT("/:categoryId")
		admin.DELETE("/:categoryId")
	}
}
