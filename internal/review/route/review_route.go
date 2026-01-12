package route

import (
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/review"
)

func ReviewRoute(route *gin.RouterGroup, reviewController review.ReviewController, notifyController review.NotificationController) {
	reviews := route.Group("/reviews")
	{
		reviews.POST("/", reviewController.CreateReview)
		reviews.PUT("/:reviewId", reviewController.UpdateReview)
		reviews.DELETE("/:reviewId", reviewController.DeleteReview)
		reviews.POST("/:reviewId/media", reviewController.UploadReviewMedia)
	}

	productReviews := route.Group("/products/:productId/reviews")
	{
		productReviews.GET("/", reviewController.ListProductReviews)
	}

	sellerReviews := route.Group("/sellers/:sellerId/reviews")
	{
		sellerReviews.GET("/", reviewController.ListSellerReviews)
	}

	userReviews := route.Group("/users/me/reviews")
	{
		userReviews.GET("/", reviewController.GetUserReviews)
	}

	admin := route.Group("/admin/reviews")
	{
		admin.GET("/pending", reviewController.ListPendingReviews)
		admin.PATCH("/:reviewId/approve", reviewController.ApproveReview)
		admin.PATCH("/:reviewId/reject", reviewController.RejectReview)
	}

	notifications := route.Group("/notifications")
	{
		notifications.GET("/", notifyController.ListNotifications)
		notifications.GET("/unread", notifyController.GetUnreadNotifications)
		notifications.PATCH("/:notificationId/read", notifyController.MarkAsRead)
		notifications.PATCH("/read-all", notifyController.MarkAllAsRead)
		notifications.DELETE("/:notificationId", notifyController.DeleteNotification)
	}
}
