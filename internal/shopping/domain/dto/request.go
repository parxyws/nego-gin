package dto

type CartItemAddRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	VariantID int32  `json:"variant_id"`
	Quantity  int32  `json:"quantity" validate:"required,min=1"`
}

type CartItemUpdateRequest struct {
	Quantity int32 `json:"quantity" validate:"required,min=1"`
}

type OrderCreateRequest struct {
	OrderType         string `json:"order_type" validate:"required"`
	AuctionID         string `json:"auction_id"`
	BillingAddressID  int32  `json:"billing_address_id" validate:"required"`
	ShippingAddressID int32  `json:"shipping_address_id" validate:"required"`
	CustomerNotes     string `json:"customer_notes"`
}

type PaymentCreateRequest struct {
	OrderID       string  `json:"order_id" validate:"required"`
	PaymentMethod string  `json:"payment_method" validate:"required"`
	Amount        float64 `json:"amount" validate:"required"`
}

type ShipmentUpdateRequest struct {
	Carrier        string `json:"carrier"`
	TrackingNumber string `json:"tracking_number"`
	ShipmentStatus string `json:"shipment_status"`
}

type CartMergeRequest struct {
	GuestSessionID string `json:"guest_session_id" validate:"required"`
}

type CheckoutValidateRequest struct {
	CartID string `json:"cart_id" validate:"required"`
}

type CheckoutCalculateRequest struct {
	CartID            string `json:"cart_id" validate:"required"`
	ShippingAddressID int32  `json:"shipping_address_id" validate:"required"`
}
