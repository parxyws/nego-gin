package dto

import "time"

type AuctionResponse struct {
	AuctionID         string    `json:"auction_id"`
	ProductID         string    `json:"product_id"`
	StartingPrice     float64   `json:"starting_price"`
	ReservePrice      float64   `json:"reserve_price"`
	BidIncrement      float64   `json:"bid_increment"`
	CurrentPrice      float64   `json:"current_price"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
	AutoExtendMinutes int32     `json:"auto_extend_minutes"`
	Status            string    `json:"status"`
	WinningBidID      string    `json:"winning_bid_id"`
	WinningBidderID   string    `json:"winning_bidder_id"`
	FinalPrice        float64   `json:"final_price"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type BidResponse struct {
	BidID        string    `json:"bid_id"`
	AuctionID    string    `json:"auction_id"`
	BidderID     string    `json:"bidder_id"`
	BidAmount    float64   `json:"bid_amount"`
	MaxBidAmount float64   `json:"max_bid_amount"`
	BidStatus    string    `json:"bid_status"`
	CreatedAt    time.Time `json:"created_at"`
}

type WatchlistResponse struct {
	WatchlistID         int32     `json:"watchlist_id"`
	UserID              string    `json:"user_id"`
	ProductID           string    `json:"product_id"`
	WatchlistType       string    `json:"watchlist_type"`
	NotifyOnOutbid      bool      `json:"notify_on_outbid"`
	NotifyBeforeEnd     bool      `json:"notify_before_end"`
	NotifyMinutesBefore int32     `json:"notify_minutes_before"`
	CreatedAt           time.Time `json:"created_at"`
}
