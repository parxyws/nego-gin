package repository

import (
	"context"

	"github.com/parxyws/nego-gin/internal/product"
	"github.com/parxyws/nego-gin/internal/product/domain"
	"gorm.io/gorm"
)

type ShopCategoryRepositoryImpl struct {
	DB *gorm.DB
}

func NewShopCategoryRepositoryImpl(db *gorm.DB) product.ShopCategoryRepository {
	return &ShopCategoryRepositoryImpl{DB: db}
}

func (s *ShopCategoryRepositoryImpl) ReadShopCategoryProductById(ctx context.Context, category *domain.ShopCategory) (*domain.ShopCategory, error) {
	foundCategory := &domain.ShopCategory{}
	tx := s.DB.WithContext(ctx)
	if err := tx.Where("shop_category_id = ?", category.ShopCategoryID).Preload("Products").First(&foundCategory).Error; err != nil {
		return nil, err
	}

	return foundCategory, nil
}

func (s *ShopCategoryRepositoryImpl) CreateShopCategory(ctx context.Context, category *domain.ShopCategory) (*domain.ShopCategory, error) {
	tx := s.DB.WithContext(ctx)

	err := tx.Transaction(func(tx *gorm.DB) error {
		if result := tx.Where("category_name = ?", category.CategoryName).FirstOrCreate(category); result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *ShopCategoryRepositoryImpl) UpdateShopCategory(ctx context.Context, category *domain.ShopCategory) (*domain.ShopCategory, error) {
	tx := s.DB.WithContext(ctx)

	err := tx.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&domain.ShopCategory{}).Where("shop_category_id = ?", category.ShopCategoryID).Updates(category)

		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *ShopCategoryRepositoryImpl) DeleteShopCategory(ctx context.Context, category *domain.ShopCategory) error {
	tx := s.DB.WithContext(ctx)

	err := tx.Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&domain.ShopCategory{}, "shop_category_id = ?", category.ShopCategoryID)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *ShopCategoryRepositoryImpl) ReadShopCategoryById(ctx context.Context, category *domain.ShopCategory) (*domain.ShopCategory, error) {
	foundCategory := &domain.ShopCategory{}
	tx := s.DB.WithContext(ctx)
	if err := tx.Where("shop_category_id = ?", category.ShopCategoryID).First(&foundCategory).Error; err != nil {
		return nil, err
	}

	return foundCategory, nil
}

func (s *ShopCategoryRepositoryImpl) ReadAllShopCategory(ctx context.Context) ([]domain.ShopCategory, error) {
	var shopCategory []domain.ShopCategory
	tx := s.DB.WithContext(ctx)
	if err := tx.Find(&shopCategory).Error; err != nil {
		return nil, err
	}

	return shopCategory, nil
}
