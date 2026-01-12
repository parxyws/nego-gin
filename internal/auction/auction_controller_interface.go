package auction

import "github.com/gin-gonic/gin"

type AuctionController interface {
	ListActiveAuctions(c *gin.Context)
	ListLiveAuctions(c *gin.Context)
	ListAuctionsEndingSoon(c *gin.Context)
	GetAuctionDetails(c *gin.Context)
	CreateAuction(c *gin.Context)
	UpdateAuction(c *gin.Context)
	CancelAuction(c *gin.Context)
	ListSellerAuctions(c *gin.Context)
}

type BidController interface {
	GetAuctionBidHistory(c *gin.Context)
	CreateAuctionBid(c *gin.Context)
	GetUserBidHistory(c *gin.Context)
	GetUserActiveBid(c *gin.Context)
	CancelAuctionBid(c *gin.Context)
}

type WatchlistController interface {
	GetUserWatchlist(c *gin.Context)
	AddUserWatchlist(c *gin.Context)
	RemoveUserWatchlist(c *gin.Context)
	UpdateUserWatchlist(c *gin.Context)
}
