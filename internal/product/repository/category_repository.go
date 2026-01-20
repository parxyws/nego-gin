package repository

import (
	"context"

	"github.com/parxyws/nego-gin/internal/product"
	"github.com/parxyws/nego-gin/internal/product/domain"
	"github.com/sirupsen/logrus"
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
			logrus.WithFields(logrus.Fields{"function": "CategoryRepository.CreateCategory", "category_name": category.CategoryName}).Warn("no rows affected during category creation")
			return result.Error
		}

		if result.Error != nil {
			logrus.WithFields(logrus.Fields{"function": "CategoryRepository.CreateCategory", "category_name": category.CategoryName}).Errorf("failed to create category: %v", result.Error)
			return result.Error
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	logrus.WithFields(logrus.Fields{"function": "CategoryRepository.CreateCategory", "category_id": category.CategoryID}).Info("category created successfully")
	return category, nil
}

func (c *CategoryRepositoryImpl) UpdateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	tx := c.DB.WithContext(ctx)

	err := tx.Transaction(func(tx *gorm.DB) error {
		result := tx.Save(&category)

		if result.Error != nil {
			logrus.WithFields(logrus.Fields{"function": "CategoryRepository.UpdateCategory", "category_id": category.CategoryID}).Errorf("failed to update category: %v", result.Error)
			return result.Error
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	logrus.WithFields(logrus.Fields{"function": "CategoryRepository.UpdateCategory", "category_id": category.CategoryID}).Info("category updated successfully")
	return category, nil
}

func (c *CategoryRepositoryImpl) DeleteCategory(ctx context.Context, category *domain.Category) error {
	tx := c.DB.WithContext(ctx)

	err := tx.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&domain.Category{}, "category_id = ?", category.CategoryID).Error; err != nil {
			logrus.WithFields(logrus.Fields{"function": "CategoryRepository.DeleteCategory", "category_id": category.CategoryID}).Errorf("failed to delete category: %v", err)
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{"function": "CategoryRepository.DeleteCategory", "category_id": category.CategoryID}).Info("category deleted successfully")
	return nil
}

func (c *CategoryRepositoryImpl) ReadCategoryById(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	foundCategory := new(domain.Category)
	tx := c.DB.WithContext(ctx)

	if err := tx.Where("category_id = ?", category.CategoryID).First(foundCategory).Error; err != nil {
		logrus.WithFields(logrus.Fields{"function": "CategoryRepository.ReadCategoryById", "category_id": category.CategoryID}).Warnf("failed to read category by id: %v", err)
		return nil, err
	}

	return foundCategory, nil
}

func (c *CategoryRepositoryImpl) ReadCategoryProductById(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	foundCategory := new(domain.Category)
	tx := c.DB.WithContext(ctx)

	if err := tx.Where("category_id = ?", category.CategoryID).Preload("products").Find(foundCategory).Error; err != nil {
		logrus.WithFields(logrus.Fields{"function": "CategoryRepository.ReadCategoryProductById", "category_id": category.CategoryID}).Warnf("failed to read category products: %v", err)
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
		logrus.WithFields(logrus.Fields{"function": "CategoryRepository.ReadAllCategory"}).Errorf("failed to read all categories: %v", err)
		return nil, err
	}

	return categories, nil
}
