package route

import (
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/product"
)

func CategoryRoute(route *gin.RouterGroup, controller product.CategoryController) {
	category := route.Group("/categories")
	{
		category.GET("/", controller.ListCategories)
		category.GET("/:categoryId", controller.GetCategory)
		category.GET("/:categoryId/products", controller.ListProductsInCategory)
	}

	admin := route.Group("/admin/categories")
	{
		admin.POST("/", controller.CreateCategory)
		admin.PUT("/:categoryId", controller.UpdateCategory)
		admin.DELETE("/:categoryId", controller.RemoveCategory)
	}
}
