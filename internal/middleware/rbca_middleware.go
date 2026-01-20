package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/user/domain"
	"github.com/parxyws/nego-gin/pkg/helper"
	"github.com/sirupsen/logrus"
)

func (m *ManagerMiddleware) RBCAMiddleware(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*JwtPayload)

		var roleList []domain.Role
		tx := m.db.WithContext(context.Background())
		if err := tx.Preload("Permissions").Where("role_name IN ?", user.Role).Find(&roleList).Error; err != nil {
			logrus.WithFields(logrus.Fields{"function": "RBCAMiddleware", "user_id": user.ID, "roles": user.Role}).Warnf("failed to fetch roles: %v", err)
			helper.Error(c, http.StatusForbidden, "access denied", err)
			c.Abort()
			return
		}

		// Check if the user has any of the required roles
		hasRequiredRole := false
		for _, requiredRole := range roles {
			for _, userRole := range user.Role {
				if requiredRole == userRole {
					hasRequiredRole = true
					break
				}
			}
			if hasRequiredRole {
				break
			}
		}

		if !hasRequiredRole {
			logrus.WithFields(logrus.Fields{"function": "RBCAMiddleware", "user_id": user.ID, "user_roles": user.Role, "required_roles": roles}).Warn("access denied: user does not have required roles")
			helper.Error(c, http.StatusForbidden, "access denied", errors.New("insufficient permissions"))
			c.Abort()
			return
		}

		c.Next()
	}
}
