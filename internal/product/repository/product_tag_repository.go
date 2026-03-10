package repository

import (
	"context"

	"github.com/parxyws/nego-gin/internal/product"
	"github.com/parxyws/nego-gin/internal/product/domain"
	"gorm.io/gorm"
)

type ProductTagRepository struct {
	DB *gorm.DB
}

func NewProductTagRepository(db *gorm.DB) product.ProductTagRepository {
	return &ProductTagRepository{DB: db}
}

func (p *ProductTagRepository) CreateProductTag(ctx context.Context, product *domain.ProductTag) (*domain.ProductTag, error) {
	tx := p.DB.WithContext(ctx)

	if result := tx.Create(product); result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

func (p *ProductTagRepository) UpdateProductTag(ctx context.Context, product *domain.ProductTag) (*domain.ProductTag, error) {
	tx := p.DB.WithContext(ctx)
	if result := tx.Model(&domain.ProductTag{}).Where("product_id = ? AND tag_id = ?", product.ProductID, product.TagID).Updates(product); result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

func (p *ProductTagRepository) DeleteProductTag(ctx context.Context, product *domain.ProductTag) error {
	tx := p.DB.WithContext(ctx)
	if result := tx.Delete(product); result.Error != nil {
		return result.Error
	}

	return nil
}

func (p *ProductTagRepository) ReadProductTagByTagId(ctx context.Context, product *domain.ProductTag) (*domain.ProductTag, error) {
	productTag := new(domain.ProductTag)
	tx := p.DB.WithContext(ctx)

	if err := tx.Where("tag_id = ?", product.TagID).First(&productTag).Error; err != nil {
		return nil, err
	}

	return productTag, nil
}

func (p *ProductTagRepository) ReadAllProductTagByProductId(ctx context.Context, product *domain.ProductTag) ([]domain.ProductTag, error) {
	var productTags []domain.ProductTag
	tx := p.DB.WithContext(ctx)

	if err := tx.Where("product_id = ?", product.ProductID).Find(productTags).Error; err != nil {
		return nil, err
	}

	return productTags, nil
}
