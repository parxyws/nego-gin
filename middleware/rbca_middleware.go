package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/user/domain"
	"github.com/parxyws/nego-gin/pkg/helper"
)

func (m *ManagerMiddleware) RBCAMiddleware(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*JwtPayload)

		var roleList []domain.Role
		tx := m.db.WithContext(context.Background())
		if err := tx.Preload("Permissions").Where("role_name IN ?", user.Role).Find(&roleList).Error; err != nil {
			helper.Error(c, http.StatusForbidden, "no access", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
