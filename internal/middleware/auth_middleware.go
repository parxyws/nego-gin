package middleware

import (
	"errors"
	"fmt"
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
	roleSlice, _ := claims["role"]

	var roles []string
	if roleSlice != nil {
		switch v := roleSlice.(type) {
		case []string:
			roles = v
		case []interface{}:
			for _, item := range v {
				if str, ok := item.(string); ok {
					roles = append(roles, str)
				}
			}
		default:
			fmt.Printf("Unexpected type for role: %T\n", v)
			roles = []string{}
		}
	} else {
		roles = []string{}
	}

	return &JwtPayload{
		ID:       claims["id"].(string),
		Username: claims["username"].(string),
		Email:    claims["email"].(string),
		Role:     roles,
	}
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

	// Extract path segment 3 — identical to how seedPermission derives resource names.
	// e.g. "/api/admin/roles" → parts[3] == "admin"
	//      "/api/sellers/products" → parts[3] == "sellers"
	//      "/api/orders/..." → parts[3] == "orders"
	parts := strings.Split(c.Request.URL.Path, "/")
	if len(parts) < 4 {
		// Route too shallow to determine a resource — allow any authenticated user.
		return true
	}
	resource := parts[3]

	switch resource {
	case "admin":
		// Protected (Admin): requires admin role.
		return slices.Contains(user.Role, "admin")
	case "sellers":
		// Protected (Seller): requires vendor role.
		return slices.Contains(user.Role, "vendor")
	default:
		// Protected: any authenticated user (customer, vendor, admin, etc.)
		return len(user.Role) > 0
	}
}

func (m *ManagerMiddleware) unauthorized(c *gin.Context, code int, message string) {
	helper.Error(c, code, message, errors.New(message))
}

//func logoutResponse() func(c *gin.Context) {
//	return func(c *gin.Context) {
//		// This demonstrates that claims are now accessible during logout
//		claims := jwt.ExtractClaims(c)
//		user, exists := c.Get(identityKey)
//
//		response := gin.H{
//			"code":    http.StatusOK,
//			"message": "Successfully logged out",
//		}
//
//		// Show that we can access user information during logout
//		if len(claims) > 0 {
//			response["logged_out_user"] = claims[identityKey]
//		}
//		if exists {
//			response["user_info"] = user.(*User).UserName
//		}
//
//		c.JSON(http.StatusOK, response)
//	}
//}
