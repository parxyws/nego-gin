package dto

type ReviewCreateRequest struct {
	OrderItemID        int32  `json:"order_item_id" validate:"required"`
	ProductID          string `json:"product_id" validate:"required"`
	SellerID           string `json:"seller_id" validate:"required"`
	ProductRating      int32  `json:"product_rating" validate:"required,min=1,max=5"`
	ProductReviewTitle string `json:"product_review_title"`
	ProductReviewText  string `json:"product_review_text"`
	SellerRating       int32  `json:"seller_rating" validate:"required,min=1,max=5"`
	SellerReviewText   string `json:"seller_review_text"`
}

type ReviewUpdateRequest struct {
	ProductRating      int32  `json:"product_rating" validate:"min=1,max=5"`
	ProductReviewTitle string `json:"product_review_title"`
	ProductReviewText  string `json:"product_review_text"`
	SellerRating       int32  `json:"seller_rating" validate:"min=1,max=5"`
	SellerReviewText   string `json:"seller_review_text"`
}

type NotificationMarkAsReadRequest struct {
	NotificationIDs []string `json:"notification_ids" validate:"required"`
}

type ReviewModerateRequest struct {
	ModerationNote string `json:"moderation_note"`
}
