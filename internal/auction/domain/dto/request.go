package dto

import "time"

type AuctionCreateRequest struct {
	ProductID         string    `json:"product_id" validate:"required"`
	StartingPrice     float64   `json:"starting_price" validate:"required"`
	ReservePrice      float64   `json:"reserve_price"`
	BidIncrement      float64   `json:"bid_increment" validate:"required"`
	StartTime         time.Time `json:"start_time" validate:"required"`
	EndTime           time.Time `json:"end_time" validate:"required"`
	AutoExtendMinutes int32     `json:"auto_extend_minutes"`
}

type AuctionUpdateRequest struct {
	StartingPrice     float64   `json:"starting_price"`
	ReservePrice      float64   `json:"reserve_price"`
	BidIncrement      float64   `json:"bid_increment"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
	AutoExtendMinutes int32     `json:"auto_extend_minutes"`
	Status            string    `json:"status"`
}

type BidCreateRequest struct {
	AuctionID    string  `json:"auction_id" validate:"required"`
	BidAmount    float64 `json:"bid_amount" validate:"required"`
	MaxBidAmount float64 `json:"max_bid_amount"`
}

type WatchlistAddRequest struct {
	ProductID           string `json:"product_id" validate:"required"`
	WatchlistType       string `json:"watchlist_type" validate:"required"`
	NotifyOnOutbid      bool   `json:"notify_on_outbid"`
	NotifyBeforeEnd     bool   `json:"notify_before_end"`
	NotifyMinutesBefore int32  `json:"notify_minutes_before"`
}
