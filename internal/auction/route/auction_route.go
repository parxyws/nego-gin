package route

import (
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/auction"
)

func AuctionRoute(route *gin.RouterGroup, controller auction.AuctionController, authMiddleware *jwt.GinJWTMiddleware) {
	// Public auction browsing
	auctions := route.Group("/auctions")
	{
		auctions.GET("", controller.ListActiveAuctions)
		auctions.GET("/live", controller.ListLiveAuctions)
		auctions.GET("/ending-soon", controller.ListAuctionsEndingSoon)
		auctions.GET("/:auctionId", controller.GetAuctionDetails)
	}

	// Seller-protected auction management
	sellerAuctions := route.Group("/sellers/auctions")
	sellerAuctions.Use(authMiddleware.MiddlewareFunc())
	{
		sellerAuctions.POST("", controller.CreateAuction)
		sellerAuctions.PUT("/:auctionId", controller.UpdateAuction)
		sellerAuctions.DELETE("/:auctionId", controller.CancelAuction)
		sellerAuctions.GET("", controller.ListSellerAuctions)
	}
}
