package repository

import (
	"context"

	"github.com/parxyws/nego-gin/internal/product"
	"github.com/parxyws/nego-gin/internal/product/domain"
	"gorm.io/gorm"
)

type ProductAttributeRepository struct {
	DB *gorm.DB
}

func NewProductAttributeRepository(db *gorm.DB) product.ProductAttributeRepository {
	return &ProductAttributeRepository{DB: db}
}

func (p *ProductAttributeRepository) CreateProductAttribute(ctx context.Context, product *domain.ProductAttribute) (*domain.ProductAttribute, error) {
	tx := p.DB.WithContext(ctx)
	if result := tx.Create(product); result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

func (p *ProductAttributeRepository) UpdateProductAttribute(ctx context.Context, product *domain.ProductAttribute) (*domain.ProductAttribute, error) {
	tx := p.DB.WithContext(ctx)
	if result := tx.Updates(product); result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

func (p *ProductAttributeRepository) DeleteProductAttribute(ctx context.Context, product *domain.ProductAttribute) error {
	tx := p.DB.WithContext(ctx)
	if result := tx.Delete(product); result.Error != nil {
		return result.Error
	}

	return nil
}

func (p *ProductAttributeRepository) ReadProductAttributeById(ctx context.Context, product *domain.ProductAttribute) (*domain.ProductAttribute, error) {
	foundProduct := new(domain.ProductAttribute)
	tx := p.DB.WithContext(ctx)
	if err := tx.Where("attribute_id = ?", product.AttributeID).First(foundProduct).Error; err != nil {
		return nil, err
	}

	return foundProduct, nil
}

func (p *ProductAttributeRepository) ReadProductAttributeByName(ctx context.Context, product *domain.ProductAttribute) (*domain.ProductAttribute, error) {
	foundProduct := new(domain.ProductAttribute)
	tx := p.DB.WithContext(ctx)
	if err := tx.Where("attribute_name = ?", product.AttributeName).First(foundProduct).Error; err != nil {
		return nil, err
	}

	return foundProduct, nil
}

func (p *ProductAttributeRepository) ReadAllProductAttributeByProductId(ctx context.Context, product *domain.ProductAttribute) ([]domain.ProductAttribute, error) {
	var productAttributes []domain.ProductAttribute
	tx := p.DB.WithContext(ctx)
	if err := tx.Where("product_id = ?", product.ProductID).Find(&productAttributes).Error; err != nil {
		return nil, err
	}

	return productAttributes, nil
}
