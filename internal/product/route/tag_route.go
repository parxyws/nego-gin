package route

import (
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/product"
)

func TagRoute(route *gin.RouterGroup, controller product.TagController, authMiddleware *jwt.GinJWTMiddleware) {
	// Public tag browsing
	tags := route.Group("/tags")
	{
		tags.GET("", controller.ListTags)
		tags.GET("/:tagId/products", controller.GetProductsByTag)
	}
}
