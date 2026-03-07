package route

import (
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/auction"
)

func BidRoute(route *gin.RouterGroup, controller auction.BidController, authMiddleware *jwt.GinJWTMiddleware) {
	auctions := route.Group("/auctions")
	{
		// Public — read bid history
		auctions.GET("/:auctionId/bids", controller.GetAuctionBidHistory)

		// Protected — place or cancel bids
		protectedAuctions := auctions.Group("/")
		protectedAuctions.Use(authMiddleware.MiddlewareFunc())
		{
			protectedAuctions.POST("/:auctionId/bids", controller.CreateAuctionBid)
			protectedAuctions.DELETE("/:auctionId/bids/:bidId", controller.CancelAuctionBid)
		}
	}

	// Protected — user's own bids
	userBids := route.Group("/users/me")
	userBids.Use(authMiddleware.MiddlewareFunc())
	{
		userBids.GET("/bids", controller.GetUserBidHistory)
		userBids.GET("/bids/active", controller.GetUserActiveBid)
	}
}
