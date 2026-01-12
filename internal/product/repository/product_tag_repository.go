package repository

import (
	"context"

	"github.com/parxyws/nego-gin/internal/product"
	"github.com/parxyws/nego-gin/internal/product/domain"
	"gorm.io/gorm"
)

type ProductTagRepositoryImpl struct {
	DB *gorm.DB
}

func NewProductTagRepositoryImpl(db *gorm.DB) product.ProductTagRepository {
	return &ProductTagRepositoryImpl{DB: db}
}

func (p *ProductTagRepositoryImpl) CreateProductTag(ctx context.Context, product *domain.ProductTag) (*domain.ProductTag, error) {
	tx := p.DB.WithContext(ctx)

	err := tx.Transaction(func(tx *gorm.DB) error {
		if result := tx.Create(product); result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *ProductTagRepositoryImpl) UpdateProductTag(ctx context.Context, product *domain.ProductTag) (*domain.ProductTag, error) {
	tx := p.DB.WithContext(ctx)
	err := tx.Transaction(func(tx *gorm.DB) error {
		if result := tx.Save(product); result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *ProductTagRepositoryImpl) DeleteProductTag(ctx context.Context, product *domain.ProductTag) error {
	tx := p.DB.WithContext(ctx)
	err := tx.Transaction(func(tx *gorm.DB) error {
		if result := tx.Delete(product); result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (p *ProductTagRepositoryImpl) ReadProductTagByTagId(ctx context.Context, product *domain.ProductTag) (*domain.ProductTag, error) {
	productTag := new(domain.ProductTag)
	tx := p.DB.WithContext(ctx)

	if err := tx.Where("tag_id = ?", product.TagID).First(&productTag).Error; err != nil {
		return nil, err
	}

	return productTag, nil
}

func (p *ProductTagRepositoryImpl) ReadAllProductTagByProductId(ctx context.Context, product *domain.ProductTag) ([]domain.ProductTag, error) {
	var productTags []domain.ProductTag
	tx := p.DB.WithContext(ctx)

	if err := tx.Where("product_id = ?", product.ProductID).Find(productTags).Error; err != nil {
		return nil, err
	}

	return productTags, nil
}
