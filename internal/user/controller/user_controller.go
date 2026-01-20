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
		helper.Error(c, http.StatusBadRequest, "Retrieve user context failed", nil)
		return
	}

	result, err := u.userService.GetCurrentUser(ctx, request.(middleware.JwtPayload))
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Retrieve user data failed", err)
		return
	}

	helper.Success(c, http.StatusOK, "User data fetched successfully", result)
}

func (u *UserControllerImpl) UpdateCurrentUser(c *gin.Context) {
	request, exists := c.Get("user")
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	if !exists {
		helper.Error(c, http.StatusBadRequest, "Retrieve user context failed", nil)
		return
	}

	body := new(dto.UpdateCurrentUserRequest)
	if err := c.ShouldBindJSON(body); err != nil {
		helper.Error(c, http.StatusBadRequest, "Parse request body failed: invalid JSON format", err)
		return
	}

	if err := validator.ValidateStruct(ctx, body); err != nil {
		helper.Error(c, http.StatusBadRequest, "Validation failed", err)
		return
	}

	result, err := u.userService.UpdateCurrentUser(ctx, request.(middleware.JwtPayload), body)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Update user profile failed", err)
		return
	}

	helper.Success(c, http.StatusOK, "User profile updated successfully", result)
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
		helper.Error(c, http.StatusBadRequest, "Retrieve user context failed", nil)
		return
	}

	if err := u.userService.DeleteCurrentUser(ctx, request.(middleware.JwtPayload)); err != nil {
		helper.Error(c, http.StatusInternalServerError, "Delete user account failed", err)
		return
	}

	helper.Success(c, http.StatusOK, "User account deleted successfully", nil)
}

func (u *UserControllerImpl) GetUserProfile(c *gin.Context) {
	request := c.Param("userId")
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	result, err := u.userService.GetUser(ctx, &dto.GetUserProfileRequest{UserID: request})
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Retrieve user profile failed", err)
		return
	}

	helper.Success(c, http.StatusOK, "User profile fetched successfully", result)
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
		helper.Error(c, http.StatusBadRequest, "Retrieve user context failed", nil)
		return
	}

	result, err := u.userAddressService.ListUserAddress(ctx, request.(middleware.JwtPayload).ID)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Retrieve address list failed", err)
		return
	}

	helper.Success(c, http.StatusOK, "Address list fetched successfully", result)
}

func (u *UserControllerImpl) CreateUserAddress(c *gin.Context) {
	request, exists := c.Get("user")
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	if !exists {
		helper.Error(c, http.StatusBadRequest, "Retrieve user context failed", nil)
		return
	}

	body := new(dto.UserAddressCreateRequest)
	body.UserId = request.(middleware.JwtPayload).ID
	if err := c.ShouldBindJSON(body); err != nil {
		helper.Error(c, http.StatusBadRequest, "Parse request body failed: invalid JSON format", err)
		return
	}

	if err := validator.ValidateStruct(ctx, body); err != nil {
		helper.Error(c, http.StatusBadRequest, "Validation failed", err)
		return
	}

	result, err := u.userAddressService.CreateUserAddress(ctx, body)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Create address failed", err)
		return
	}

	helper.Success(c, http.StatusOK, "Address created successfully", result)
}

func (u *UserControllerImpl) UpdateUserAddress(c *gin.Context) {
	request, exists := c.Get("user")
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	if !exists {
		helper.Error(c, http.StatusBadRequest, "Retrieve user context failed", nil)
		return
	}

	body := new(dto.UserAddressUpdateRequest)
	body.UserId = request.(middleware.JwtPayload).ID
	if err := c.ShouldBindJSON(body); err != nil {
		helper.Error(c, http.StatusBadRequest, "Parse request body failed: invalid JSON format", err)
		return
	}

	result, err := u.userAddressService.UpdateUserAddress(ctx, body)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Update address failed", err)
		return
	}

	helper.Success(c, http.StatusOK, "Address updated successfully", result)
}

func (u *UserControllerImpl) DeleteUserAddress(c *gin.Context) {
	request, exists := c.Get("user")
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	if !exists {
		helper.Error(c, http.StatusBadRequest, "Retrieve user context failed", nil)
		return
	}

	params, err := strconv.ParseInt(c.Param("addressId"), 10, 32)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "Validation failed: invalid address identifier", err)
		return
	}

	if err := u.userAddressService.DeleteUserAddress(ctx, int32(params), request.(middleware.JwtPayload).ID); err != nil {
		helper.Error(c, http.StatusInternalServerError, "Delete address failed", err)
		return
	}

	helper.Success(c, http.StatusOK, "Address deleted successfully", nil)
}

func (u *UserControllerImpl) SetDefaultAddress(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}
