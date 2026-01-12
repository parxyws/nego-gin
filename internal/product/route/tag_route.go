package route

import (
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/product"
)

func TagRoute(route *gin.RouterGroup, controller product.TagController) {
	tags := route.Group("/tags")
	{
		tags.GET("/")
		tags.GET("/:tag_id/products")
	}
}
