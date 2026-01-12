package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/user"
	"github.com/parxyws/nego-gin/internal/user/domain/dto"
	"github.com/parxyws/nego-gin/middleware"
	"github.com/parxyws/nego-gin/pkg/helper"
	"github.com/parxyws/nego-gin/pkg/validator"
)

type UserControllerImpl struct {
	userService user.UserService
}

func NewUserControllerImpl(userService user.UserService) user.UserController {
	return &UserControllerImpl{userService: userService}
}

func (u *UserControllerImpl) GetCurrentUser(c *gin.Context) {
	request, exists := c.Get("user")
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	if !exists {
		helper.Error(c, http.StatusBadRequest, "failed to get current user", nil)
	}

	result, err := u.userService.GetCurrentUser(ctx, request.(middleware.JwtPayload))
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "failed to get current user", nil)
	}

	helper.Success(c, http.StatusOK, "success to get user", result)
}

func (u *UserControllerImpl) UpdateCurrentUser(c *gin.Context) {
	request, exists := c.Get("user")
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	if !exists {
		helper.Error(c, http.StatusBadRequest, "failed to get current user", nil)
	}

	body := new(dto.UpdateCurrentUserRequest)
	if err := c.ShouldBindJSON(body); err != nil {
		helper.Error(c, http.StatusBadRequest, "failed to parse body", nil)
	}

	if err := validator.ValidateStruct(ctx, body); err != nil {
		helper.Error(c, http.StatusBadRequest, "failed to validate body", nil)
	}

	result, err := u.userService.UpdateCurrentUser(ctx, request.(middleware.JwtPayload), body)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "failed to update current user", nil)
	}

	helper.Success(c, http.StatusOK, "success to update current user", result)
}

func (u *UserControllerImpl) UpdateAvatar(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (u *UserControllerImpl) DeleteCurrentUser(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (u *UserControllerImpl) GetUserProfile(c *gin.Context) {
	request := c.Params.ByName("userId")
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	result, err := u.userService.GetUser(ctx, &dto.GetUserProfileRequest{UserID: request})
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "failed to get user profile", nil)
	}

	helper.Success(c, http.StatusOK, "success to get user profile", result)
}

func (u *UserControllerImpl) GetUserSellerRatings(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (u *UserControllerImpl) GetUserListOfAddresses(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (u *UserControllerImpl) CreateUserAddress(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (u *UserControllerImpl) UpdateUserAddress(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (u *UserControllerImpl) DeleteUserAddress(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (u *UserControllerImpl) SetDefaultAddress(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}
