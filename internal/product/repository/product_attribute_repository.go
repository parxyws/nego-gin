package repository

import (
	"context"

	"github.com/parxyws/nego-gin/internal/product"
	"github.com/parxyws/nego-gin/internal/product/domain"
	"gorm.io/gorm"
)

type ProductAttributeRepositoryImpl struct {
	DB *gorm.DB
}

func NewProductAttributeRepositoryImpl(db *gorm.DB) product.ProductAttributeRepository {
	return &ProductAttributeRepositoryImpl{DB: db}
}

func (p *ProductAttributeRepositoryImpl) CreateProductAttribute(ctx context.Context, product *domain.ProductAttribute) (*domain.ProductAttribute, error) {
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

func (p *ProductAttributeRepositoryImpl) UpdateProductAttribute(ctx context.Context, product *domain.ProductAttribute) (*domain.ProductAttribute, error) {
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

func (p *ProductAttributeRepositoryImpl) DeleteProductAttribute(ctx context.Context, product *domain.ProductAttribute) error {
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

func (p *ProductAttributeRepositoryImpl) ReadProductAttributeById(ctx context.Context, product *domain.ProductAttribute) (*domain.ProductAttribute, error) {
	foundProduct := new(domain.ProductAttribute)
	tx := p.DB.WithContext(ctx)
	if err := tx.Where("attribute_id = ?", product.AttributeID).First(foundProduct).Error; err != nil {
		return nil, err
	}

	return foundProduct, nil
}

func (p *ProductAttributeRepositoryImpl) ReadProductAttributeByName(ctx context.Context, product *domain.ProductAttribute) (*domain.ProductAttribute, error) {
	foundProduct := new(domain.ProductAttribute)
	tx := p.DB.WithContext(ctx)
	if err := tx.Where("attribute_name = ?", product.AttributeName).First(foundProduct).Error; err != nil {
		return nil, err
	}

	return foundProduct, nil
}

func (p *ProductAttributeRepositoryImpl) ReadAllProductAttributeByProductId(ctx context.Context, product *domain.ProductAttribute) ([]domain.ProductAttribute, error) {
	var productAttributes []domain.ProductAttribute
	tx := p.DB.WithContext(ctx)
	if err := tx.Where("product_id = ?", product.ProductID).Find(&productAttributes).Error; err != nil {
		return nil, err
	}

	return productAttributes, nil
}
