package repository

import (
	"context"

	"github.com/parxyws/nego-gin/internal/product"
	"github.com/parxyws/nego-gin/internal/product/domain"
	"gorm.io/gorm"
)

type ProductRepository struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) product.ProductRepository {
	return &ProductRepository{DB: db}
}

func (p *ProductRepository) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	tx := p.DB.WithContext(ctx)

	if result := tx.Create(product); result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

func (p *ProductRepository) UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	tx := p.DB.WithContext(ctx)

	if result := tx.Save(product); result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

func (p *ProductRepository) DeleteProduct(ctx context.Context, product *domain.Product) error {
	tx := p.DB.WithContext(ctx)

	if result := tx.Delete(product); result.Error != nil {
		return result.Error
	}

	return nil
}

func (p *ProductRepository) ReadProductById(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	foundProduct := new(domain.Product)
	tx := p.DB.WithContext(ctx)
	if err := tx.Where("product_id = ?", product.ProductID).First(&foundProduct).Error; err != nil {
		return nil, err
	}

	return foundProduct, nil
}

func (p *ProductRepository) ReadAllProduct(ctx context.Context, product *domain.Product) ([]domain.Product, error) {
	var products []domain.Product

	tx := p.DB.WithContext(ctx)
	if err := tx.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (p *ProductRepository) ReadAllProductSlug(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	foundProduct := new(domain.Product)
	tx := p.DB.WithContext(ctx)
	if err := tx.Where("slug = ?", product.Slug).First(&foundProduct).Error; err != nil {
		return nil, err
	}

	return foundProduct, nil
}
