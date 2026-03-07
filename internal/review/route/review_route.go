package route

import (
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/review"
)

func ReviewRoute(route *gin.RouterGroup, reviewController review.ReviewController, notifyController review.NotificationController, authMiddleware *jwt.GinJWTMiddleware) {
	// Public — read product and seller reviews
	route.GET("/products/:productId/reviews", reviewController.ListProductReviews)
	route.GET("/sellers/:sellerId/reviews", reviewController.ListSellerReviews)

	// Protected — create/update/delete reviews and view own reviews
	reviews := route.Group("/reviews")
	reviews.Use(authMiddleware.MiddlewareFunc())
	{
		reviews.POST("", reviewController.CreateReview)
		reviews.PUT("/:reviewId", reviewController.UpdateReview)
		reviews.DELETE("/:reviewId", reviewController.DeleteReview)
		reviews.POST("/:reviewId/media", reviewController.UploadReviewMedia)
	}

	// Protected — user's own reviews
	userReviews := route.Group("/users/me/reviews")
	userReviews.Use(authMiddleware.MiddlewareFunc())
	{
		userReviews.GET("", reviewController.GetUserReviews)
	}

	// Admin — review moderation
	adminReviews := route.Group("/admin/reviews")
	adminReviews.Use(authMiddleware.MiddlewareFunc())
	{
		adminReviews.GET("/pending", reviewController.ListPendingReviews)
		adminReviews.PATCH("/:reviewId/approve", reviewController.ApproveReview)
		adminReviews.PATCH("/:reviewId/reject", reviewController.RejectReview)
	}

	// Protected — notifications
	notifications := route.Group("/notifications")
	notifications.Use(authMiddleware.MiddlewareFunc())
	{
		notifications.GET("", notifyController.ListNotifications)
		notifications.GET("/unread", notifyController.GetUnreadNotifications)
		notifications.PATCH("/:notificationId/read", notifyController.MarkAsRead)
		notifications.PATCH("/read-all", notifyController.MarkAllAsRead)
		notifications.DELETE("/:notificationId", notifyController.DeleteNotification)
	}
}
