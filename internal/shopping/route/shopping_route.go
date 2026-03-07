package route

import (
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/shopping"
)

func ShoppingRoute(
	route *gin.RouterGroup,
	cartController shopping.CartController,
	orderController shopping.OrderController,
	paymentController shopping.PaymentController,
	shipmentController shopping.ShipmentController,
	authMiddleware *jwt.GinJWTMiddleware,
) {
	// Cart — Public/Protected (supports both guest and authenticated per README)
	cart := route.Group("/cart")
	{
		cart.GET("", cartController.GetCart)
		cart.POST("/items", cartController.AddItemToCart)
		cart.PUT("/items/:cartItemId", cartController.UpdateCartItem)
		cart.DELETE("/items/:cartItemId", cartController.RemoveItemFromCart)
		cart.DELETE("", cartController.ClearCart)

		// Merge cart requires authentication (after login)
		protectedCart := cart.Group("/")
		protectedCart.Use(authMiddleware.MiddlewareFunc())
		{
			protectedCart.POST("/merge", cartController.MergeCart)
		}
	}

	// Checkout — Public/Protected
	checkout := route.Group("/checkout")
	{
		checkout.POST("/validate", orderController.ValidateCheckout)
		checkout.POST("/calculate", orderController.CalculateCheckout)

		// Creating an order requires an account
		protectedCheckout := checkout.Group("/")
		protectedCheckout.Use(authMiddleware.MiddlewareFunc())
		{
			protectedCheckout.POST("", orderController.CreateOrder)
		}
	}

	// Orders — Protected
	orders := route.Group("/orders")
	orders.Use(authMiddleware.MiddlewareFunc())
	{
		orders.GET("", orderController.ListUserOrders)
		orders.GET("/:orderId", orderController.GetOrderDetails)
		orders.POST("/:orderId/cancel", orderController.CancelOrder)
		orders.GET("/:orderId/items", orderController.ListOrderItems)
		orders.GET("/:orderId/payments", paymentController.GetOrderPayments)
		orders.GET("/:orderId/shipments", shipmentController.GetOrderShipments)
	}

	// Payments — Protected
	payments := route.Group("/payments")
	payments.Use(authMiddleware.MiddlewareFunc())
	{
		payments.POST("", paymentController.CreatePayment)
		payments.GET("/:paymentId", paymentController.GetPaymentDetails)
		payments.POST("/:paymentId/confirm", paymentController.ConfirmPayment)
		payments.POST("/:paymentId/refund", paymentController.RefundPayment)
	}

	// Shipments — Protected reads / Public tracking number
	shipments := route.Group("/shipments")
	{
		shipments.GET("/:shipmentId/tracking", shipmentController.GetTrackingHistory) // Public

		protectedShipments := shipments.Group("/")
		protectedShipments.Use(authMiddleware.MiddlewareFunc())
		{
			protectedShipments.GET("/:shipmentId", shipmentController.GetShipmentDetails)
		}
	}

	// Seller-protected order and shipment management
	seller := route.Group("/sellers")
	seller.Use(authMiddleware.MiddlewareFunc())
	{
		sellerOrders := seller.Group("/orders")
		{
			sellerOrders.GET("", orderController.ListSellerOrders)
			sellerOrders.PATCH("/:orderId/status", orderController.UpdateOrderStatus)
			sellerOrders.POST("/:orderId/shipments", shipmentController.CreateShipment)
		}

		sellerShipments := seller.Group("/shipments")
		{
			sellerShipments.PUT("/:shipmentId", shipmentController.UpdateShipment)
		}
	}

	// Payment webhooks — Public (verified by signature in handler)
	webhooks := route.Group("/webhooks")
	{
		webhooks.POST("/payment/:provider", paymentController.PaymentWebhook)
	}
}
