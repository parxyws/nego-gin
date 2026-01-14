package controller

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/middleware"
	"github.com/parxyws/nego-gin/internal/user"
	"github.com/parxyws/nego-gin/internal/user/domain/dto"
	"github.com/parxyws/nego-gin/pkg/helper"
	"github.com/parxyws/nego-gin/pkg/validator"
)

type UserControllerImpl struct {
	userService        user.UserService
	userAddressService user.UserAddressService
}

func NewUserController(userService user.UserService, userAddressService user.UserAddressService) user.UserController {
	return &UserControllerImpl{userService: userService, userAddressService: userAddressService}
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
	request, exists := c.Get("user")
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	if !exists {
		helper.Error(c, http.StatusBadRequest, "failed to get current user", nil)
	}

	if err := u.userService.DeleteCurrentUser(ctx, request.(middleware.JwtPayload)); err != nil {
		helper.Error(c, http.StatusInternalServerError, "failed to delete current user", nil)
	}

	helper.Success(c, http.StatusOK, "success to delete current user", nil)
}

func (u *UserControllerImpl) GetUserProfile(c *gin.Context) {
	request := c.Param("userId")
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

func (u *UserControllerImpl) ListUserAddress(c *gin.Context) {
	request, exists := c.Get("user")
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	if !exists {
		helper.Error(c, http.StatusBadRequest, "failed to get current user", nil)
	}

	result, err := u.userAddressService.ListUserAddress(ctx, request.(middleware.JwtPayload).ID)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "failed to list user address", nil)
	}

	helper.Success(c, http.StatusOK, "success to list user address", result)
}

func (u *UserControllerImpl) CreateUserAddress(c *gin.Context) {
	request, exists := c.Get("user")
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	if !exists {
		helper.Error(c, http.StatusBadRequest, "failed to get current user", nil)
	}

	body := new(dto.UserAddressCreateRequest)
	body.UserId = request.(middleware.JwtPayload).ID
	if err := c.ShouldBindJSON(body); err != nil {
		helper.Error(c, http.StatusBadRequest, "failed to parse body", nil)
	}

	if err := validator.ValidateStruct(ctx, body); err != nil {
		helper.Error(c, http.StatusBadRequest, "failed to validate body", nil)
	}

	result, err := u.userAddressService.CreateUserAddress(ctx, body)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "failed to create user address", nil)
	}

	helper.Success(c, http.StatusOK, "success to create user address", result)
}

func (u *UserControllerImpl) UpdateUserAddress(c *gin.Context) {
	request, exists := c.Get("user")
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	if !exists {
		helper.Error(c, http.StatusBadRequest, "failed to get current user", nil)
	}

	body := new(dto.UserAddressUpdateRequest)
	body.UserId = request.(middleware.JwtPayload).ID
	if err := c.ShouldBindJSON(body); err != nil {
		helper.Error(c, http.StatusBadRequest, "failed to parse body", nil)
	}

	result, err := u.userAddressService.UpdateUserAddress(ctx, body)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "failed to update user address", nil)
	}

	helper.Success(c, http.StatusOK, "success to update user address", result)
}

func (u *UserControllerImpl) DeleteUserAddress(c *gin.Context) {
	request, exists := c.Get("user")
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	if !exists {
		helper.Error(c, http.StatusBadRequest, "failed to get current user", nil)
	}

	params, err := strconv.ParseInt(c.Param("addressId"), 10, 32)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "failed to parse addressId", nil)
	}

	if err := u.userAddressService.DeleteUserAddress(ctx, int32(params), request.(middleware.JwtPayload).ID); err != nil {
		helper.Error(c, http.StatusInternalServerError, "failed to delete user address", nil)
	}

	helper.Success(c, http.StatusOK, "success to delete user address", nil)
}

func (u *UserControllerImpl) SetDefaultAddress(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}
