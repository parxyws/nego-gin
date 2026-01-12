package route

import (
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/auction"
)

func AuctionRoute(route *gin.RouterGroup, controller auction.AuctionController) {
	auctions := route.Group("/auctions")
	{
		auctions.GET("/", controller.ListActiveAuctions)
		auctions.GET("/live", controller.ListLiveAuctions)
		auctions.GET("/ending-soon", controller.ListAuctionsEndingSoon)
		auctions.GET("/:auctionId", controller.GetAuctionDetails)
	}

	seller := route.Group("/sellers/auctions")
	{
		seller.POST("/", controller.CreateAuction)
		seller.PUT("/:auctionId", controller.UpdateAuction)
		seller.DELETE("/:auctionId", controller.CancelAuction)
		seller.GET("/", controller.ListSellerAuctions)
	}
}
