package middleware

import (
	"errors"
	"fmt"
	"slices"
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
		Timeout:         time.Hour,
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
