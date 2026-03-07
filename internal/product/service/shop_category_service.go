package service

import (
	"context"

	"github.com/parxyws/nego-gin/internal/product"
	"github.com/parxyws/nego-gin/internal/product/domain/dto"
)

type ShopCategoryService struct {
	shopCtRepo product.ShopCategoryService
}

func NewShopCategoryService(shopCtRepo product.ShopCategoryService) product.ShopCategoryService {
	return &ShopCategoryService{shopCtRepo: shopCtRepo}
}

func (s *ShopCategoryService) ListShopCategories(ctx context.Context, sellerId string) ([]dto.ShopCategoryResponse, error) {
	panic("implement me")
}

func (s *ShopCategoryService) GetShopCategory(ctx context.Context, categoryId int32) (*dto.ShopCategoryResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ShopCategoryService) CreateShopCategory(ctx context.Context, entity *dto.ShopCategoryCreateRequest) (*dto.ShopCategoryResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ShopCategoryService) UpdateShopCategory(ctx context.Context, entity *dto.ShopCategoryUpdateRequest) (*dto.ShopCategoryResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ShopCategoryService) RemoveShopCategory(ctx context.Context, categoryId int32) error {
	//TODO implement me
	panic("implement me")
}
