package middleware

import (
	"errors"
	"net/http"
	"slices"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/parxyws/nego-gin/pkg/helper"
)

type JwtPayload struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Role     []string `json:"role"`
}

func (m *ManagerMiddleware) JWTAuthMiddleware() (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(m.initJwtParams())
	if err != nil {
		return nil, err
	}

	return authMiddleware, nil
}

func (m *ManagerMiddleware) initJwtParams() *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:           "nego API",
		Key:             []byte(m.cfg.Server.JWTSecretKey),
		Timeout:         time.Minute * 15,
		MaxRefresh:      time.Hour,
		IdentityKey:     "user",
		PayloadFunc:     m.payloadHandler,
		IdentityHandler: m.identityHandler,
		Authorizer:      m.authorizer,
		Unauthorized:    m.unauthorized,
		LogoutResponse:  m.logoutResponse,
		TokenLookup:     "header: Authorization",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	}
}

func (m *ManagerMiddleware) payloadHandler(data any) jwt2.MapClaims {
	if v, ok := data.(*JwtPayload); ok {
		return jwt2.MapClaims{
			"id":       v.ID,
			"username": v.Username,
			"email":    v.Email,
			"role":     v.Role,
		}
	}

	return jwt2.MapClaims{}
}

func (m *ManagerMiddleware) identityHandler(c *gin.Context) any {
	claims := jwt.ExtractClaims(c)

	payload := &JwtPayload{}

	if id, ok := claims["id"].(string); ok {
		payload.ID = id
	}
	if username, ok := claims["username"].(string); ok {
		payload.Username = username
	}
	if email, ok := claims["email"].(string); ok {
		payload.Email = email
	}

	if roleSlice, ok := claims["role"]; ok && roleSlice != nil {
		switch v := roleSlice.(type) {
		case []string:
			payload.Role = v
		case []any:
			var roles []string
			for _, item := range v {
				if str, ok := item.(string); ok {
					roles = append(roles, str)
				}
			}
			payload.Role = roles
		}
	}

	return payload
}

// authorizer implements the JWT-layer RBAC gate.
// It mirrors the same resource-extraction logic used by the seeder:
// the "resource" is URL path segment 3 (index in split "/" path).
//
// Policy rules (derived from seed + README):
//   - "superadmin" → passes all routes unconditionally
//   - resource == "admin"   → requires "admin" role
//   - resource == "sellers" → requires "vendor" role
//   - any other resource    → any authenticated user passes
func (m *ManagerMiddleware) authorizer(c *gin.Context, data any) bool {
	user, ok := data.(*JwtPayload)
	if !ok {
		return false
	}

	// Superadmin bypasses all role checks.
	if slices.Contains(user.Role, "superadmin") {
		return true
	}

	path := c.Request.URL.Path

	// Protected (Admin): requires admin role.
	if strings.HasPrefix(path, "/api/v1/admin/") || path == "/api/v1/admin" {
		return slices.Contains(user.Role, "admin")
	}

	// Protected (Seller): requires vendor role.
	if strings.HasPrefix(path, "/api/v1/sellers/") || path == "/api/v1/sellers" {
		return slices.Contains(user.Role, "vendor")
	}

	// Protected: any authenticated user (customer, vendor, admin, etc.)
	return len(user.Role) > 0
}

func (m *ManagerMiddleware) unauthorized(c *gin.Context, code int, message string) {
	helper.Error(c, code, message, errors.New(message))
}

func (m *ManagerMiddleware) logoutResponse(c *gin.Context) {
	// This demonstrates that claims are accessible during logout
	claims := jwt.ExtractClaims(c)
	var username string
	if val, ok := claims["username"].(string); ok {
		username = val
	}

	// // Clear the refresh token cookie
	c.SetCookie("refresh_token", "", -1, "/api/auth/refresh", "", true, true)

	helper.Success(c, http.StatusOK, "Successfully logged out", gin.H{
		"logged_out_user": username,
	})
}
