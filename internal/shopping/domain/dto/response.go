package dto

import "time"

type CartResponse struct {
	CartID    string             `json:"cart_id"`
	UserID    string             `json:"user_id"`
	SessionID string             `json:"session_id"`
	Items     []CartItemResponse `json:"items,omitempty"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type CartItemResponse struct {
	CartItemID    int32     `json:"cart_item_id"`
	ProductID     string    `json:"product_id"`
	VariantID     int32     `json:"variant_id"`
	Quantity      int32     `json:"quantity"`
	PriceSnapshot float64   `json:"price_snapshot"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type OrderResponse struct {
	OrderID        string              `json:"order_id"`
	UserID         string              `json:"user_id"`
	OrderNumber    string              `json:"order_number"`
	OrderType      string              `json:"order_type"`
	Subtotal       float64             `json:"subtotal"`
	TaxAmount      float64             `json:"tax_amount"`
	ShippingAmount float64             `json:"shipping_amount"`
	DiscountAmount float64             `json:"discount_amount"`
	TotalAmount    float64             `json:"total_amount"`
	OrderStatus    string              `json:"order_status"`
	Items          []OrderItemResponse `json:"items,omitempty"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
}

type OrderItemResponse struct {
	OrderItemID int32     `json:"order_item_id"`
	ProductID   string    `json:"product_id"`
	VariantID   int32     `json:"variant_id"`
	ProductName string    `json:"product_name"`
	Sku         string    `json:"sku"`
	Quantity    int32     `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	Subtotal    float64   `json:"subtotal"`
	CreatedAt   time.Time `json:"created_at"`
}

type PaymentResponse struct {
	PaymentID     string    `json:"payment_id"`
	OrderID       string    `json:"order_id"`
	PaymentMethod string    `json:"payment_method"`
	TransactionID string    `json:"transaction_id"`
	Amount        float64   `json:"amount"`
	CurrencyCode  string    `json:"currency_code"`
	PaymentStatus string    `json:"payment_status"`
	ProcessedAt   time.Time `json:"processed_at"`
	CreatedAt     time.Time `json:"created_at"`
}

type ShipmentResponse struct {
	ShipmentID        int32             `json:"shipment_id"`
	OrderID           string            `json:"order_id"`
	Carrier           string            `json:"carrier"`
	TrackingNumber    string            `json:"tracking_number"`
	ShipmentStatus    string            `json:"shipment_status"`
	EstimatedDelivery time.Time         `json:"estimated_delivery"`
	ActualDelivery    time.Time         `json:"actual_delivery"`
	TrackingHistory   []TrackingHistory `json:"tracking_history,omitempty"`
}

type TrackingHistory struct {
	Status      string    `json:"status"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
	TrackedAt   time.Time `json:"tracked_at"`
}

type CheckoutValidateResponse struct {
	IsValid bool     `json:"is_valid"`
	Errors  []string `json:"errors,omitempty"`
}

type CheckoutCalculateResponse struct {
	Subtotal       float64 `json:"subtotal"`
	TaxAmount      float64 `json:"tax_amount"`
	ShippingAmount float64 `json:"shipping_amount"`
	TotalAmount    float64 `json:"total_amount"`
}
