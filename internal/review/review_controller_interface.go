package review

import "github.com/gin-gonic/gin"

type ReviewController interface {
	ListProductReviews(c *gin.Context)
	ListSellerReviews(c *gin.Context)
	CreateReview(c *gin.Context)
	UpdateReview(c *gin.Context)
	DeleteReview(c *gin.Context)
	GetUserReviews(c *gin.Context)
	UploadReviewMedia(c *gin.Context)

	// Admin Moderation
	ListPendingReviews(c *gin.Context)
	ApproveReview(c *gin.Context)
	RejectReview(c *gin.Context)
}

type NotificationController interface {
	ListNotifications(c *gin.Context)
	GetUnreadNotifications(c *gin.Context)
	MarkAsRead(c *gin.Context)
	MarkAllAsRead(c *gin.Context)
	DeleteNotification(c *gin.Context)
}
