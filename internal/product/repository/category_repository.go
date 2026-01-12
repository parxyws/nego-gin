package repository

import (
	"context"

	"github.com/parxyws/nego-gin/internal/product"
	"github.com/parxyws/nego-gin/internal/product/domain"
	"gorm.io/gorm"
)

type CategoryRepositoryImpl struct {
	DB *gorm.DB
}

func NewCategoryRepositoryImpl(db *gorm.DB) product.CategoryRepository {
	return &CategoryRepositoryImpl{DB: db}
}

func (c *CategoryRepositoryImpl) CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	tx := c.DB.WithContext(ctx)

	err := tx.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(&category)

		if result.RowsAffected == 0 {
			return result.Error
		}

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

func (c *CategoryRepositoryImpl) UpdateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	tx := c.DB.WithContext(ctx)

	err := tx.Transaction(func(tx *gorm.DB) error {
		result := tx.Save(&category)

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

func (c *CategoryRepositoryImpl) DeleteCategory(ctx context.Context, category *domain.Category) error {
	tx := c.DB.WithContext(ctx)

	err := tx.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&domain.Category{}, "category_id = ?", category.CategoryID).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (c *CategoryRepositoryImpl) ReadCategoryById(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	foundCategory := new(domain.Category)
	tx := c.DB.WithContext(ctx)

	if err := tx.Where("category_id = ?", category.CategoryID).First(foundCategory).Error; err != nil {
		return nil, err
	}

	return foundCategory, nil
}

func (c *CategoryRepositoryImpl) ReadCategoryByName(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	foundCategory := new(domain.Category)
	tx := c.DB.WithContext(ctx)

	if err := tx.Where("category_name = ?", category.CategoryName).First(foundCategory).Error; err != nil {
		return nil, err
	}

	return foundCategory, nil
}

func (c *CategoryRepositoryImpl) ReadAllCategory(ctx context.Context) ([]domain.Category, error) {
	var categories []domain.Category
	tx := c.DB.WithContext(ctx)

	if err := tx.Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}
