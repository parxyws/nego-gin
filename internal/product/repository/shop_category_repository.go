package repository

import (
	"context"
	"fmt"

	"github.com/parxyws/nego-gin/internal/product"
	"github.com/parxyws/nego-gin/internal/product/domain"
	"gorm.io/gorm"
)

type ShopCategoryRepository struct {
	DB *gorm.DB
}

func NewShopCategoryRepository(db *gorm.DB) product.ShopCategoryRepository {
	return &ShopCategoryRepository{DB: db}
}

func (s *ShopCategoryRepository) CreateShopCategory(ctx context.Context, category *domain.ShopCategory) (*domain.ShopCategory, error) {
	tx := s.DB.WithContext(ctx)

	result := tx.Where("category_name = ?", category.CategoryName).FirstOrCreate(category)
	if result.Error != nil {
		return nil, result.Error
	}

	return category, nil
}

func (s *ShopCategoryRepository) UpdateShopCategory(ctx context.Context, category *domain.ShopCategory) (*domain.ShopCategory, error) {
	tx := s.DB.WithContext(ctx)

	if result := tx.Model(&domain.ShopCategory{}).Where("shop_category_id = ?", category.ShopCategoryID).Updates(category); result.Error != nil {
		return nil, result.Error
	}

	return category, nil
}

func (s *ShopCategoryRepository) DeleteShopCategory(ctx context.Context, category *domain.ShopCategory) error {
	tx := s.DB.WithContext(ctx)

	if result := tx.Delete(&domain.ShopCategory{}, "shop_category_id = ?", category.ShopCategoryID); result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *ShopCategoryRepository) ReadShopCategoryProductById(ctx context.Context, category *domain.ShopCategory) (*domain.ShopCategory, error) {
	foundCategory := &domain.ShopCategory{}
	tx := s.DB.WithContext(ctx)

	if err := tx.Where("shop_category_id = ?", category.ShopCategoryID).Preload("Products").First(&foundCategory).Error; err != nil {
		return nil, fmt.Errorf("ShopCategoryRepository.ReadShopCategoryProductById - %w", err)
	}

	return foundCategory, nil
}

func (s *ShopCategoryRepository) ReadShopCategoryById(ctx context.Context, category *domain.ShopCategory) (*domain.ShopCategory, error) {
	foundCategory := &domain.ShopCategory{}
	tx := s.DB.WithContext(ctx)

	if err := tx.Where("shop_category_id = ?", category.ShopCategoryID).First(&foundCategory).Error; err != nil {
		return nil, err
	}

	return foundCategory, nil
}

func (s *ShopCategoryRepository) ReadAllShopCategory(ctx context.Context) ([]domain.ShopCategory, error) {
	var shopCategories []domain.ShopCategory
	tx := s.DB.WithContext(ctx)

	if err := tx.Find(&shopCategories).Error; err != nil {
		return nil, err
	}

	return shopCategories, nil
}
