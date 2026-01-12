package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/pkg/helper"

	"github.com/parxyws/nego-gin/internal/user"
	"github.com/parxyws/nego-gin/internal/user/domain/dto"
	"github.com/parxyws/nego-gin/middleware"
	"github.com/parxyws/nego-gin/pkg/validator"
)

type AuthControllerImpl struct {
	authService user.AuthService
}

func NewAuthController(authService user.AuthService) user.AuthController {
	return &AuthControllerImpl{authService: authService}
}

func (a *AuthControllerImpl) Logout(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthControllerImpl) VerifyEmail(c *gin.Context) {
	req := new(dto.UserValidateAccRequest)
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	if err := c.ShouldBindJSON(req); err != nil {
		helper.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := validator.ValidateStruct(ctx, req); err != nil {
		helper.Error(c, http.StatusBadRequest, "Validation failed", err)
		return
	}

	if err := a.authService.ValidateUser(ctx, req); err != nil {
		helper.Error(c, http.StatusUnauthorized, "Validation User failed", err)
		return
	}

	helper.Success(c, http.StatusOK, "Validate success", nil)
}

func (a *AuthControllerImpl) Register(c *gin.Context) {
	req := new(dto.UserRegisterRequest)
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	if err := c.ShouldBindJSON(req); err != nil {
		helper.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := validator.ValidateStruct(ctx, req); err != nil {
		helper.Error(c, http.StatusBadRequest, "Validation failed", err)
		return
	}

	result, err := a.authService.Register(ctx, req)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Register failed", err)
		return
	}

	helper.Success(c, http.StatusOK, "Register success", result)
}

func (a *AuthControllerImpl) Login(c *gin.Context) {
	request := new(dto.UserLoginRequest)
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	if err := c.ShouldBindJSON(request); err != nil {
		helper.Error(c, http.StatusBadRequest, "Invalid request body", err)
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

	helper.Success(c, http.StatusOK, "Login success", result)
}

func (a *AuthControllerImpl) RefreshToken(c *gin.Context) {
	token := new(dto.JwtToken)
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	if err := c.ShouldBindJSON(token); err != nil {
		helper.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	value, exists := c.Get("user")
	if !exists {
		helper.Error(c, http.StatusUnauthorized, "Token not found", nil)
		return
	}

	result, err := a.authService.RefreshToken(ctx, token, value.(*middleware.JwtPayload))
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Refresh token failed", err)
		return
	}

	helper.Success(c, http.StatusOK, "Refresh token success", result)
}

func (a *AuthControllerImpl) ForgotPassword(c *gin.Context) {
	req := new(dto.ForgotPasswordRequest)
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	if err := c.ShouldBindJSON(req); err != nil {
		helper.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := validator.ValidateStruct(ctx, req); err != nil {
		helper.Error(c, http.StatusBadRequest, "Validation failed", err)
		return
	}

	if err := a.authService.ForgotPassword(ctx, req); err != nil {
		helper.Error(c, http.StatusInternalServerError, "Forgot password failed", err)
		return
	}

	helper.Success(c, http.StatusOK, "OTP sent to your email", nil)
}

func (a *AuthControllerImpl) ResetPassword(c *gin.Context) {
	req := new(dto.ResetPasswordRequest)
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	if err := c.ShouldBindJSON(req); err != nil {
		helper.Error(c, http.StatusBadRequest, "Invalid request body", err)
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

	helper.Success(c, http.StatusOK, "Password reset success", nil)
}

func (a *AuthControllerImpl) ResendVerification(c *gin.Context) {
	req := new(dto.UserRegisterResponse)
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	if err := c.ShouldBindJSON(req); err != nil {
		helper.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := validator.ValidateStruct(ctx, req); err != nil {
		helper.Error(c, http.StatusBadRequest, "Validation failed", err)
		return
	}

	err := a.authService.ResendVerification(ctx, req)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Resend verification failed", err)
		return
	}

	helper.Success(c, http.StatusOK, "Resend verification success", nil)
}
