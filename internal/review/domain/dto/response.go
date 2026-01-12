package dto

import "time"

type ReviewResponse struct {
	ReviewID           string                `json:"review_id"`
	ProductID          string                `json:"product_id"`
	ReviewerID         string                `json:"reviewer_id"`
	ProductRating      int32                 `json:"product_rating"`
	ProductReviewTitle string                `json:"product_review_title"`
	ProductReviewText  string                `json:"product_review_text"`
	SellerRating       int32                 `json:"seller_rating"`
	SellerReviewText   string                `json:"seller_review_text"`
	Media              []ReviewMediaResponse `json:"media,omitempty"`
	CreatedAt          time.Time             `json:"created_at"`
}

type ReviewMediaResponse struct {
	MediaID   int32     `json:"media_id"`
	MediaType string    `json:"media_type"`
	MediaURL  string    `json:"media_url"`
	CreatedAt time.Time `json:"created_at"`
}

type NotificationResponse struct {
	NotificationID   string    `json:"notification_id"`
	NotificationType string    `json:"notification_type"`
	Title            string    `json:"title"`
	Message          string    `json:"message"`
	RelatedEntityID  string    `json:"related_entity_id"`
	ActionURL        string    `json:"action_url"`
	IsRead           bool      `json:"is_read"`
	CreatedAt        time.Time `json:"created_at"`
}
