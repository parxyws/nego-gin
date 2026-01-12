package route

import (
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/product"
)

func ShopCategoryRoute(route *gin.RouterGroup, controller product.ShopCategoryController) {
	shopCategories := route.Group("/sellers/shop-categories")
	{
		shopCategories.GET("/", controller.ListShopCategories)
		shopCategories.GET("/:shopCategoryId", controller.GetShopCategory)
		shopCategories.POST("/", controller.CreateShopCategory)
		shopCategories.PUT("/:shopCategoryId", controller.UpdateShopCategory)
		shopCategories.DELETE("/:shopCategoryId", controller.DeleteShopCategory)
	}

	public := route.Group("/shops/:sellerId/categories")
	{
		public.GET("/", controller.ListShopCategories)
		public.GET("/:shopCategoryId", controller.GetShopCategory)
	}
}
