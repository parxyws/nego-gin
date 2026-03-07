package route

import (
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/auction"
)

func WatchlistRoute(route *gin.RouterGroup, controller auction.WatchlistController, authMiddleware *jwt.GinJWTMiddleware) {
	// All watchlist operations are protected — require authentication
	watchlist := route.Group("/users/me/watchlist")
	watchlist.Use(authMiddleware.MiddlewareFunc())
	{
		watchlist.GET("", controller.GetUserWatchlist)
		watchlist.POST("", controller.AddUserWatchlist)
		watchlist.DELETE("/:watchlistId", controller.RemoveUserWatchlist)
		watchlist.PUT("/:watchlistId", controller.UpdateUserWatchlist)
	}
}
