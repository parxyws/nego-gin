package middleware

import (
	"slices"
	"time"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/parxyws/nego-gin/internal/user/domain"
)

type JwtPayload struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Role     []string `json:"role"`
}

func (m *ManagerMiddleware) InitJwtParams() *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:           "ego",
		Key:             []byte("secret"),
		Timeout:         time.Hour,
		MaxRefresh:      time.Hour,
		IdentityKey:     "id",
		PayloadFunc:     m.payloadHandler,
		IdentityHandler: m.identityHandler,
		Authorizer:      m.authorizer,
		Unauthorized:    m.unauthorized,
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
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
	role, _ := claims["role"]
	return &domain.User{
		UserID:    claims["id"].(string),
		Username:  claims["username"].(string),
		Email:     claims["email"].(string),
		UserRoles: role.([]domain.UserRole),
	}
}

func (m *ManagerMiddleware) authorizer(c *gin.Context, data any) bool {
	user, ok := data.(*JwtPayload)
	if !ok {
		return false
	}

	if slices.Contains(user.Role, "admin") {
		return true
	}

	return false
}

func (m *ManagerMiddleware) unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
		"path":    c.Request.URL.Path,
		"method":  c.Request.Method,
	})
	return
}
