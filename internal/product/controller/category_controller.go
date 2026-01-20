package controller

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/product"
	"github.com/parxyws/nego-gin/internal/product/domain/dto"
	"github.com/parxyws/nego-gin/pkg/helper"
	"github.com/parxyws/nego-gin/pkg/validator"
)

type CategoryControllerImpl struct {
	catService product.CategoryService
}

func NewCategoryControllerImpl(catService product.CategoryService) product.CategoryController {
	return &CategoryControllerImpl{catService: catService}
}

func (cat *CategoryControllerImpl) ListCategories(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	result, err := cat.catService.ListCategories(ctx)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Retrieve category list failed", err)
		return
	}

	helper.Success(c, http.StatusOK, "Category list fetched successfully", result)
}

func (cat *CategoryControllerImpl) GetCategory(c *gin.Context) {
	request := c.Params.ByName("categoryId")
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	convertReq, err := strconv.ParseInt(request, 10, 32)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "Validation failed: invalid category identifier", err)
		return
	}

	result, err := cat.catService.GetCategory(ctx, int32(convertReq))

	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Retrieve category failed", err)
		return
	}

	helper.Success(c, http.StatusOK, "Category fetched successfully", result)
}

func (cat *CategoryControllerImpl) ListProductsInCategory(c *gin.Context) {
	request := c.Params.ByName("categoryId")
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	convertReq, err := strconv.ParseInt(request, 10, 32)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "Validation failed: invalid category identifier", err)
		return
	}

	result, err := cat.catService.ListProductsInCategory(ctx, int32(convertReq))

	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Retrieve products in category failed", err)
		return
	}

	helper.Success(c, http.StatusOK, "Product list for category fetched successfully", result)

}

func (cat *CategoryControllerImpl) CreateCategory(c *gin.Context) {
	request := new(dto.CategoryCreateRequest)
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	if err := c.ShouldBindJSON(request); err != nil {
		helper.Error(c, http.StatusBadRequest, "Parse request body failed: invalid JSON format", err)
		return
	}

	if err := validator.ValidateStruct(ctx, request); err != nil {
		helper.Error(c, http.StatusBadRequest, "Validation failed", err)
		return
	}

	result, err := cat.catService.CreateCategory(ctx, request)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Create category failed", err)
		return
	}

	helper.Success(c, http.StatusOK, "Category created successfully", result)
}

func (cat *CategoryControllerImpl) UpdateCategory(c *gin.Context) {
	request := new(dto.CategoryUpdateRequest)
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	if err := c.ShouldBindJSON(request); err != nil {
		helper.Error(c, http.StatusBadRequest, "Parse request body failed: invalid JSON format", err)
		return
	}

	if err := validator.ValidateStruct(ctx, request); err != nil {
		helper.Error(c, http.StatusBadRequest, "Validation failed", err)
		return
	}

	result, err := cat.catService.UpdateCategory(ctx, request)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Update category failed", err)
		return
	}

	helper.Success(c, http.StatusOK, "Category updated successfully", result)
}

func (cat *CategoryControllerImpl) RemoveCategory(c *gin.Context) {
	request := c.Params.ByName("categoryId")
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute)
	defer cancel()

	convertReq, err := strconv.ParseInt(request, 10, 32)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "Validation failed: invalid category identifier", err)
		return
	}

	if err := cat.catService.RemoveCategory(ctx, int32(convertReq)); err != nil {
		helper.Error(c, http.StatusInternalServerError, "Delete category failed", err)
		return
	}

	helper.Success(c, http.StatusOK, "Category deleted successfully", nil)
}
