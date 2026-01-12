package shopping

import "github.com/gin-gonic/gin"

type CartController interface {
	GetCart(c *gin.Context)
	AddItemToCart(c *gin.Context)
	UpdateCartItem(c *gin.Context)
	RemoveItemFromCart(c *gin.Context)
	ClearCart(c *gin.Context)
	MergeCart(c *gin.Context)
}

type OrderController interface {
	ValidateCheckout(c *gin.Context)
	CreateOrder(c *gin.Context)
	CalculateCheckout(c *gin.Context)
	ListUserOrders(c *gin.Context)
	GetOrderDetails(c *gin.Context)
	CancelOrder(c *gin.Context)
	ListSellerOrders(c *gin.Context)
	UpdateOrderStatus(c *gin.Context)
	ListOrderItems(c *gin.Context)
}

type PaymentController interface {
	CreatePayment(c *gin.Context)
	GetPaymentDetails(c *gin.Context)
	ConfirmPayment(c *gin.Context)
	RefundPayment(c *gin.Context)
	GetOrderPayments(c *gin.Context)
	PaymentWebhook(c *gin.Context)
}

type ShipmentController interface {
	GetOrderShipments(c *gin.Context)
	GetShipmentDetails(c *gin.Context)
	GetTrackingHistory(c *gin.Context)
	CreateShipment(c *gin.Context)
	UpdateShipment(c *gin.Context)
}
