package controller

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/middleware"
	"github.com/parxyws/nego-gin/pkg/helper"

	"github.com/parxyws/nego-gin/internal/user"
	"github.com/parxyws/nego-gin/internal/user/domain/dto"
	"github.com/parxyws/nego-gin/pkg/validator"
)

type AuthController struct {
	authService user.AuthService
}

func NewAuthController(authService user.AuthService) user.AuthController {
	return &AuthController{authService: authService}
}

func (a *AuthController) Logout(c *gin.Context) {
	ctx, cancel := helper.GetContext(c)
	defer cancel()

	// Extract Bearer token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 8 {
		helper.Error(c, http.StatusUnauthorized, "Authentication failed: missing token", nil)
		return
	}
	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	if err := a.authService.Logout(ctx, accessToken); err != nil {
		helper.Error(c, http.StatusInternalServerError, "Logout failed", err)
		return
	}

	// Clear the refresh token cookie
	c.SetCookie("refresh_token", "", -1, "/api/auth/refresh", "", true, true)

	helper.Success(c, http.StatusOK, "Logged out successfully", nil)
}

func (a *AuthController) VerifyEmail(c *gin.Context) {
	req := new(dto.UserValidateAccRequest)
	ctx, cancel := helper.GetContext(c)
	defer cancel()

	if err := c.ShouldBindJSON(req); err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid request payload", err)
		return
	}

	if err := validator.ValidateStruct(ctx, req); err != nil {
		helper.Error(c, http.StatusBadRequest, "validation failed", err)
		return
	}

	if err := a.authService.ValidateUser(ctx, req); err != nil {
		helper.Error(c, http.StatusUnauthorized, "User verification failed", err)
		return
	}

	helper.Success(c, http.StatusOK, "User verified successfully", nil)
}

func (a *AuthController) Register(c *gin.Context) {
	req := new(dto.UserRegisterRequest)
	ctx, cancel := helper.GetContext(c)
	defer cancel()

	if err := c.ShouldBindJSON(req); err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid request payload", err)
		return
	}

	if err := validator.ValidateStruct(ctx, req); err != nil {
		helper.Error(c, http.StatusBadRequest, "validation failed", err)
		return
	}

	result, err := a.authService.Register(ctx, req)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Register user failed", err)
		return
	}

	helper.Success(c, http.StatusOK, "User registered successfully", result)
}

func (a *AuthController) Login(c *gin.Context) {
	request := new(dto.UserLoginRequest)
	ctx, cancel := helper.GetContext(c)
	defer cancel()

	if err := c.ShouldBindJSON(request); err != nil {
		helper.Error(c, http.StatusBadRequest, "Parse request body failed: invalid JSON format", err)
		return
	}

	if err := validator.ValidateStruct(ctx, request); err != nil {
		helper.Error(c, http.StatusBadRequest, "Validation failed", err)
		return
	}

	result, err := a.authService.Login(ctx, request)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Login failed", err)
		return
	}

	// Set refresh token as HTTP-only cookie

	c.SetCookie(
		"refresh_token",                      // Cookie name
		result.Jwt.RefreshToken,              // Refresh token value
		int(time.Hour*24*7)/int(time.Second), // 7 days expiry
		"/api/auth/refresh",                  // Only accessible to refresh endpoint
		"",                                   // Current domain
		true,                                 // Secure in production
		true,                                 // HttpOnly (cannot be accessed by JavaScript)
	)

	// Return only access token in response (for Authorization header)
	//loginResponse := struct {
	//	AccessToken string      `json:"access_token"`
	//	ExpiresIn   int64       `json:"expires_in"`
	//	TokenType   string      `json:"token_type"`
	//	User        interface{} `json:"user"`
	//}{
	//	AccessToken: result.Jwt.AccessToken,
	//	ExpiresIn:   result.Jwt.ExpiresIn,
	//	TokenType:   "Bearer",
	//	User:        result.User,
	//}

	helper.Success(c, http.StatusOK, "Login successfully", result)
}

func (a *AuthController) RefreshToken(c *gin.Context) {
	token := new(dto.JwtToken)
	ctx, cancel := helper.GetContext(c)
	defer cancel()

	if err := c.ShouldBindJSON(token); err != nil {
		helper.Error(c, http.StatusBadRequest, "Parse request body failed: invalid JSON format", err)
		return
	}

	value, exists := c.Get("user")
	if !exists {
		helper.Error(c, http.StatusUnauthorized, "Authentication failed: token not found", nil)
		return
	}

	result, err := a.authService.RefreshToken(ctx, token, value.(*middleware.JwtPayload))
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Refresh token failed", err)
		return
	}

	helper.Success(c, http.StatusOK, "Token refreshed successfully", result)
}

func (a *AuthController) ForgotPassword(c *gin.Context) {
	req := new(dto.ForgotPasswordRequest)
	ctx, cancel := helper.GetContext(c)
	defer cancel()

	if err := c.ShouldBindJSON(req); err != nil {
		helper.Error(c, http.StatusBadRequest, "Parse request body failed: invalid JSON format", err)
		return
	}

	if err := validator.ValidateStruct(ctx, req); err != nil {
		helper.Error(c, http.StatusBadRequest, "Validation failed", err)
		return
	}

	if err := a.authService.ForgotPassword(ctx, req); err != nil {
		helper.Error(c, http.StatusInternalServerError, "Process forgot password failed", err)
		return
	}

	helper.Success(c, http.StatusOK, "OTP sent successfully", nil)
}

func (a *AuthController) ResetPassword(c *gin.Context) {
	req := new(dto.ResetPasswordRequest)
	ctx, cancel := helper.GetContext(c)
	defer cancel()

	if err := c.ShouldBindJSON(req); err != nil {
		helper.Error(c, http.StatusBadRequest, "Parse request body failed: invalid JSON format", err)
		return
	}

	if err := validator.ValidateStruct(ctx, req); err != nil {
		helper.Error(c, http.StatusBadRequest, "Validation failed", err)
		return
	}

	if err := a.authService.ResetPassword(ctx, req); err != nil {
		helper.Error(c, http.StatusInternalServerError, "Reset password failed", err)
		return
	}

	helper.Success(c, http.StatusOK, "Password reset successfully", nil)
}

func (a *AuthController) ResendVerification(c *gin.Context) {
	req := new(dto.ResendVerificationRequest)
	ctx, cancel := helper.GetContext(c)
	defer cancel()

	if err := c.ShouldBindJSON(req); err != nil {
		helper.Error(c, http.StatusBadRequest, "Parse request body failed: invalid JSON format", err)
		return
	}

	if err := validator.ValidateStruct(ctx, req); err != nil {
		helper.Error(c, http.StatusBadRequest, "Validation failed", err)
		return
	}

	err := a.authService.ResendVerification(ctx, req)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Resend verification OTP failed", err)
		return
	}

	helper.Success(c, http.StatusOK, "Verification OTP resent successfully", nil)
}
