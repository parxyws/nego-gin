package repository

import (
	"context"

	"github.com/parxyws/nego-gin/internal/product"
	"github.com/parxyws/nego-gin/internal/product/domain"
	"gorm.io/gorm"
)

type ProductVariantRepositoryImpl struct {
	DB *gorm.DB
}

func NewProductVariantRepositoryImpl(db *gorm.DB) product.ProductVariantRepository {
	return &ProductVariantRepositoryImpl{DB: db}
}

func (p *ProductVariantRepositoryImpl) CreateProductVariant(ctx context.Context, product *domain.ProductVariant) (*domain.ProductVariant, error) {
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

func (p *ProductVariantRepositoryImpl) UpdateProductVariant(ctx context.Context, product *domain.ProductVariant) (*domain.ProductVariant, error) {
	tx := p.DB.WithContext(ctx)
	err := tx.Transaction(func(tx *gorm.DB) error {
		if result := tx.Updates(product); result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *ProductVariantRepositoryImpl) DeleteProductVariant(ctx context.Context, product *domain.ProductVariant) error {
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

func (p *ProductVariantRepositoryImpl) ReadProductVariantById(ctx context.Context, product *domain.ProductVariant) (*domain.ProductVariant, error) {
	foundProductVariant := new(domain.ProductVariant)
	tx := p.DB.WithContext(ctx)
	if err := tx.Where("variant_id = ?", product.VariantID).First(foundProductVariant).Error; err != nil {
		return nil, err
	}

	return foundProductVariant, nil
}

func (p *ProductVariantRepositoryImpl) ReadProductVariantByName(ctx context.Context, product *domain.ProductVariant) (*domain.ProductVariant, error) {
	foundProductVariant := new(domain.ProductVariant)
	tx := p.DB.WithContext(ctx)
	if err := tx.Where("variant_name = ?", product.VariantName).First(foundProductVariant).Error; err != nil {
		return nil, err
	}

	return foundProductVariant, nil
}

func (p *ProductVariantRepositoryImpl) ReadAllProductVariantByProductId(ctx context.Context, product *domain.ProductVariant) ([]domain.ProductVariant, error) {
	var productVariant []domain.ProductVariant

	tx := p.DB.WithContext(ctx)
	if err := tx.Where("product_id = ?", product.ProductID).Find(&productVariant).Error; err != nil {
		return nil, err
	}

	return productVariant, nil
}
