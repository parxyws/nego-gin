package route

import (
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/product"
)

func ProductVariantRoute(route *gin.RouterGroup, controller product.ProductVariantsController) {
	products := route.Group("/products")
	{
		products.GET("/:productId/variants", controller.ListProductVariants)
	}

	seller := route.Group("/sellers/products")
	{
		seller.POST("/:productId/variants", controller.CreateProductVariant)
		seller.PUT("/:productId/variants/:variantId", controller.UpdateProductVariant)
		seller.DELETE("/:productId/variants/:variantId", controller.DeleteProductVariant)
	}
}
