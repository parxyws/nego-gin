package route

import (
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/shopping"
)

func ShoppingRoute(route *gin.RouterGroup, cartController shopping.CartController, orderController shopping.OrderController, paymentController shopping.PaymentController, shipmentController shopping.ShipmentController) {
	cart := route.Group("/cart")
	{
		cart.GET("/", cartController.GetCart)
		cart.POST("/items", cartController.AddItemToCart)
		cart.PUT("/items/:cartItemId", cartController.UpdateCartItem)
		cart.DELETE("/items/:cartItemId", cartController.RemoveItemFromCart)
		cart.DELETE("/", cartController.ClearCart)
		cart.POST("/merge", cartController.MergeCart)
	}

	checkout := route.Group("/checkout")
	{
		checkout.POST("/validate", orderController.ValidateCheckout)
		checkout.POST("/", orderController.CreateOrder)
		checkout.POST("/calculate", orderController.CalculateCheckout)
	}

	orders := route.Group("/orders")
	{
		orders.GET("/", orderController.ListUserOrders)
		orders.GET("/:orderId", orderController.GetOrderDetails)
		orders.POST("/:orderId/cancel", orderController.CancelOrder)
		orders.GET("/:orderId/items", orderController.ListOrderItems)
		orders.GET("/:orderId/payments", paymentController.GetOrderPayments)
		orders.GET("/:orderId/shipments", shipmentController.GetOrderShipments)
	}

	payments := route.Group("/payments")
	{
		payments.POST("/", paymentController.CreatePayment)
		payments.GET("/:paymentId", paymentController.GetPaymentDetails)
		payments.POST("/:paymentId/confirm", paymentController.ConfirmPayment)
		payments.POST("/:paymentId/refund", paymentController.RefundPayment)
	}

	shipments := route.Group("/shipments")
	{
		shipments.GET("/:shipmentId", shipmentController.GetShipmentDetails)
		shipments.GET("/:shipmentId/tracking", shipmentController.GetTrackingHistory)
	}

	seller := route.Group("/sellers")
	{
		sellerOrders := seller.Group("/orders")
		{
			sellerOrders.GET("/", orderController.ListSellerOrders)
			sellerOrders.PATCH("/:orderId/status", orderController.UpdateOrderStatus)
			sellerOrders.POST("/:orderId/shipments", shipmentController.CreateShipment)
		}

		sellerShipments := seller.Group("/shipments")
		{
			sellerShipments.PUT("/:shipmentId", shipmentController.UpdateShipment)
		}
	}

	webhooks := route.Group("/webhooks")
	{
		webhooks.POST("/payment/:provider", paymentController.PaymentWebhook)
	}
}
