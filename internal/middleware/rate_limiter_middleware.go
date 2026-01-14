package middleware

import "github.com/gin-gonic/gin"

func (m *ManagerMiddleware) RateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
