package repository

import (
	"context"
	"fmt"

	"github.com/parxyws/nego-gin/internal/product"
	"github.com/parxyws/nego-gin/internal/product/domain"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	DB *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) product.CategoryRepository {
	return &CategoryRepository{DB: db}
}

func (c *CategoryRepository) CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	tx := c.DB.WithContext(ctx)

	result := tx.Create(&category)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("CategoryRepository.CreateCategory: failed to insert category, 0 rows affected")
	}

	return category, nil
}

func (c *CategoryRepository) UpdateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	tx := c.DB.WithContext(ctx)

	result := tx.Model(&domain.Category{}).Where("category_id = ?", category.CategoryID).Updates(category)

	if result.Error != nil {
		return nil, result.Error
	}

	return category, nil
}

func (c *CategoryRepository) DeleteCategory(ctx context.Context, category *domain.Category) error {
	tx := c.DB.WithContext(ctx)

	if err := tx.Delete(&domain.Category{}, "category_id = ?", category.CategoryID).Error; err != nil {
		return err
	}

	return nil
}

func (c *CategoryRepository) ReadCategoryById(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	foundCategory := new(domain.Category)
	tx := c.DB.WithContext(ctx)

	if err := tx.Where("category_id = ?", category.CategoryID).First(foundCategory).Error; err != nil {
		return nil, err
	}

	return foundCategory, nil
}

func (c *CategoryRepository) ReadCategoryProductById(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	foundCategory := new(domain.Category)
	tx := c.DB.WithContext(ctx)

	if err := tx.Where("category_id = ?", category.CategoryID).Preload("Products").First(foundCategory).Error; err != nil {
		return nil, fmt.Errorf("CategoryRepository.ReadCategoryProductById - %w", err)
	}

	return foundCategory, nil
}

func (c *CategoryRepository) ReadCategoryByName(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	foundCategory := new(domain.Category)
	tx := c.DB.WithContext(ctx)

	if err := tx.Where("category_name = ?", category.CategoryName).First(foundCategory).Error; err != nil {
		return nil, err
	}

	return foundCategory, nil
}

func (c *CategoryRepository) ReadAllCategory(ctx context.Context) ([]domain.Category, error) {
	var categories []domain.Category
	tx := c.DB.WithContext(ctx)

	if err := tx.Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}
