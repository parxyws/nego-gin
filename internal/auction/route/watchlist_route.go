package route

import (
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/auction"
)

func WatchlistRoute(route *gin.RouterGroup, controller auction.WatchlistController) {
	watchlist := route.Group("/users/me/watchlist")
	{
		watchlist.GET("/", controller.GetUserWatchlist)
		watchlist.POST("/", controller.AddUserWatchlist)
		watchlist.DELETE("/:watchlistId", controller.RemoveUserWatchlist)
		watchlist.PUT("/:watchlistId", controller.UpdateUserWatchlist)
	}
}
