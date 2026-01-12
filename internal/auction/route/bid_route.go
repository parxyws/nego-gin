package route

import (
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/auction"
)

func BidRoute(route *gin.RouterGroup, controller auction.BidController) {
	auctions := route.Group("/auctions")
	{
		auctions.GET("/:auctionId/bids", controller.GetAuctionBidHistory)
		auctions.POST("/:auctionId/bids", controller.CreateAuctionBid)
		auctions.DELETE("/:auctionId/bids/:bidId", controller.CancelAuctionBid)
	}

	users := route.Group("/users/me")
	{
		users.GET("/bids", controller.GetUserBidHistory)
		users.GET("/bids/active", controller.GetUserActiveBid)
	}
}
