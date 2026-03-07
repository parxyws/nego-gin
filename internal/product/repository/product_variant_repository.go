package repository

import (
	"context"

	"github.com/parxyws/nego-gin/internal/product"
	"github.com/parxyws/nego-gin/internal/product/domain"
	"gorm.io/gorm"
)

type ProductVariantRepository struct {
	DB *gorm.DB
}

func NewProductVariantRepository(db *gorm.DB) product.ProductVariantRepository {
	return &ProductVariantRepository{DB: db}
}

func (p *ProductVariantRepository) CreateProductVariant(ctx context.Context, product *domain.ProductVariant) (*domain.ProductVariant, error) {
	tx := p.DB.WithContext(ctx)

	if result := tx.Create(product); result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

func (p *ProductVariantRepository) UpdateProductVariant(ctx context.Context, product *domain.ProductVariant) (*domain.ProductVariant, error) {
	tx := p.DB.WithContext(ctx)
	if result := tx.Updates(product); result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

func (p *ProductVariantRepository) DeleteProductVariant(ctx context.Context, product *domain.ProductVariant) error {
	tx := p.DB.WithContext(ctx)
	if result := tx.Delete(product); result.Error != nil {
		return result.Error
	}

	return nil
}

func (p *ProductVariantRepository) ReadProductVariantById(ctx context.Context, product *domain.ProductVariant) (*domain.ProductVariant, error) {
	foundProductVariant := new(domain.ProductVariant)
	tx := p.DB.WithContext(ctx)
	if err := tx.Where("variant_id = ?", product.VariantID).First(foundProductVariant).Error; err != nil {
		return nil, err
	}

	return foundProductVariant, nil
}

func (p *ProductVariantRepository) ReadProductVariantByName(ctx context.Context, product *domain.ProductVariant) (*domain.ProductVariant, error) {
	foundProductVariant := new(domain.ProductVariant)
	tx := p.DB.WithContext(ctx)
	if err := tx.Where("variant_name = ?", product.VariantName).First(foundProductVariant).Error; err != nil {
		return nil, err
	}

	return foundProductVariant, nil
}

func (p *ProductVariantRepository) ReadAllProductVariantByProductId(ctx context.Context, product *domain.ProductVariant) ([]domain.ProductVariant, error) {
	var productVariant []domain.ProductVariant

	tx := p.DB.WithContext(ctx)
	if err := tx.Where("product_id = ?", product.ProductID).Find(&productVariant).Error; err != nil {
		return nil, err
	}

	return productVariant, nil
}
