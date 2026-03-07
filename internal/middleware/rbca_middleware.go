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

// LoadPermissionCache fetches every role and its associated permissions from
// the database once and stores them in an in-memory map:
//
//	role name → set of permission Names (e.g. "admin:read")
//
// Call this during application startup (after seeding) so that
// RequirePermission never has to hit the database on a hot path.
func (m *ManagerMiddleware) LoadPermissionCache() error {
	var roles []domain.Role
	if err := m.db.WithContext(context.Background()).
		Preload("Permissions").
		Find(&roles).Error; err != nil {
		return err
	}

	cache := make(map[string]map[string]bool, len(roles))
	for _, role := range roles {
		perms := make(map[string]bool, len(role.Permissions))
		for _, p := range role.Permissions {
			perms[p.Name] = true
		}
		cache[role.RoleName] = perms
	}

	m.mu.Lock()
	m.permissionCache = cache
	m.mu.Unlock()

	m.logger.WithField("roles_cached", len(cache)).Info("permission cache loaded")
	return nil
}

// ReloadPermissionCache invalidates and rebuilds the permission cache.
// Call this after any admin operation that modifies role-permission assignments.
func (m *ManagerMiddleware) ReloadPermissionCache() error {
	return m.LoadPermissionCache()
}

// RBCAMiddleware is a role-level gate: it passes only if the authenticated
// user holds at least one of the specified role names. No DB call is made.
//
// Usage:
//
//	router.Use(authMiddleware.MiddlewareFunc(), m.RBCAMiddleware("admin", "superadmin"))
func (m *ManagerMiddleware) RBCAMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*JwtPayload)

		log := logrus.WithFields(logrus.Fields{
			"function":       "RBCAMiddleware",
			"user_id":        user.ID,
			"user_roles":     user.Role,
			"required_roles": roles,
		})

		for _, requiredRole := range roles {
			for _, userRole := range user.Role {
				if requiredRole == userRole {
					c.Next()
					return
				}
			}
		}

		log.Warn("access denied: user does not have required role")
		helper.Error(c, http.StatusForbidden, "access denied", errors.New("insufficient role"))
		c.Abort()
	}
}

// RequirePermission is a permission-level gate. It checks the in-memory
// permission cache (populated at startup by LoadPermissionCache) — zero DB
// calls on the hot path.
//
// It passes only if the user's roles collectively cover ALL of the specified
// permission names (format: "resource:action", e.g. "admin:read").
//
// Superadmin bypasses this check unconditionally.
//
// Usage (single permission):
//
//	roles.GET("", m.RequirePermission("admin:read"), controller.GetListOfRoles)
//
// Usage (multiple — user must have ALL):
//
//	roles.POST("", m.RequirePermission("admin:create"), controller.CreateRole)
func (m *ManagerMiddleware) RequirePermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*JwtPayload)

		log := logrus.WithFields(logrus.Fields{
			"function":             "RequirePermission",
			"user_id":              user.ID,
			"user_roles":           user.Role,
			"required_permissions": permissions,
		})

		// Superadmin bypasses all permission checks.
		for _, r := range user.Role {
			if r == "superadmin" {
				c.Next()
				return
			}
		}

		// Read from the in-memory cache (no DB call).
		m.mu.RLock()
		cache := m.permissionCache
		m.mu.RUnlock()

		// Build the union of all permissions the user holds across their roles.
		userPerms := make(map[string]bool)
		for _, roleName := range user.Role {
			if rolePerms, ok := cache[roleName]; ok {
				for perm := range rolePerms {
					userPerms[perm] = true
				}
			}
		}

		// Check every required permission (AND logic).
		for _, required := range permissions {
			if !userPerms[required] {
				log.Warnf("access denied: missing permission %q", required)
				helper.Error(c, http.StatusForbidden, "access denied", errors.New("insufficient permissions"))
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
