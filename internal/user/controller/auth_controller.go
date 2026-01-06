package controller

import (
	"context"
	"go/types"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-fiber/pkg/validator"
	api "github.com/parxyws/nego-gin/internal/shared/dto"
	"github.com/parxyws/nego-gin/internal/user"
	"github.com/parxyws/nego-gin/internal/user/domain/dto"
)

type AuthControllerImpl struct {
	authService user.AuthService
}

func NewAuthController(authService user.AuthService) user.AuthController {
	return &AuthControllerImpl{authService: authService}
}

func (a *AuthControllerImpl) Register(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthControllerImpl) ValidateUser(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthControllerImpl) RefreshToken(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthControllerImpl) ForgotPassword(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthControllerImpl) ResetPassword(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthControllerImpl) ResendVerification(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthControllerImpl) RegisterNewUser(c *gin.Context) {
	req := new(dto.UserRegisterRequest)

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, api.ApiResponse[dto.UserRegisterResponse]{
			Success: false,
			Message: "Invalid request body",
			Error:   err,
		})
		return
	}

	if err := validator.ValidateStruct(ctx, req); err != nil {
		c.JSON(http.StatusBadRequest, api.ApiResponse[dto.UserRegisterResponse]{
			Success: false,
			Message: "Validation failed",
			Error:   err,
		})
		return
	}

	result, err := a.authService.Register(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ApiResponse[dto.UserRegisterResponse]{
			Success: false,
			Message: "Register failed",
			Error:   err,
		})
		return
	}

	c.JSON(http.StatusOK, api.ApiResponse[dto.UserRegisterResponse]{
		Success: true,
		Message: "Register success",
		Data:    *result,
	})
}

func (a *AuthControllerImpl) ValidateNewUser(c *gin.Context) {
	request := new(dto.UserValidateAccRequest)
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(http.StatusBadRequest, api.ApiResponse[types.Nil]{
			Success: false,
			Message: "Invalid request body",
			Error:   err,
		})
		return
	}

	if err := validator.ValidateStruct(ctx, request); err != nil {
		c.JSON(http.StatusBadRequest, api.ApiResponse[types.Nil]{
			Success: false,
			Message: "Validation failed",
			Error:   err,
		})
		return
	}

	if err := a.authService.ValidateUser(ctx, request); err != nil {
		c.JSON(http.StatusUnauthorized, api.ApiResponse[types.Nil]{
			Success: false,
			Message: "Validation failed",
			Error:   err,
		})
		return
	}

	c.JSON(http.StatusOK, api.ApiResponse[types.Nil]{
		Success: true,
		Message: "Validate success",
	})
}

func (a *AuthControllerImpl) Login(c *gin.Context) {
	request := new(dto.UserLoginRequest)
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(http.StatusBadRequest, api.ApiResponse[dto.UserResponse]{
			Success: false,
			Message: "Invalid request body",
			Error:   err,
		})
		return
	}

	if err := validator.ValidateStruct(ctx, request); err != nil {
		c.JSON(http.StatusBadRequest, api.ApiResponse[dto.UserResponse]{
			Success: false,
			Message: "Validation failed",
			Error:   err,
		})
		return
	}

	result, err := a.authService.Login(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ApiResponse[dto.UserResponse]{
			Success: false,
			Message: "Login failed",
			Error:   err,
		})
		return
	}

	c.JSON(http.StatusOK, api.ApiResponse[dto.UserResponse]{
		Success: true,
		Message: "Login success",
		Data:    *result,
	})

}
