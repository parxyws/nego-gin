package service

import (
	"context"

	"github.com/parxyws/nego-gin/internal/product"
	"github.com/parxyws/nego-gin/internal/product/domain/dto"
)

type ShopCategoryServiceImpl struct {
	shopCtRepo product.ShopCategoryService
}

func NewShopCategoryServiceImpl(shopCtRepo product.ShopCategoryService) product.ShopCategoryService {
	return &ShopCategoryServiceImpl{shopCtRepo: shopCtRepo}
}

func (s *ShopCategoryServiceImpl) ListShopCategories(ctx context.Context, sellerId string) ([]dto.ShopCategoryResponse, error) {
	panic("implement me")
}

func (s *ShopCategoryServiceImpl) GetShopCategory(ctx context.Context, categoryId int32) (*dto.ShopCategoryResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ShopCategoryServiceImpl) CreateShopCategory(ctx context.Context, entity *dto.ShopCategoryCreateRequest) (*dto.ShopCategoryResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ShopCategoryServiceImpl) UpdateShopCategory(ctx context.Context, entity *dto.ShopCategoryUpdateRequest) (*dto.ShopCategoryResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ShopCategoryServiceImpl) RemoveShopCategory(ctx context.Context, categoryId int32) error {
	//TODO implement me
	panic("implement me")
}
