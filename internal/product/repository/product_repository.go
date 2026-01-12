package repository

import (
	"context"

	"github.com/parxyws/nego-gin/internal/product"
	"github.com/parxyws/nego-gin/internal/product/domain"
	"gorm.io/gorm"
)

type ProductRepositoryImpl struct {
	DB *gorm.DB
}

func NewProductRepositoryImpl(db *gorm.DB) product.ProductRepository {
	return &ProductRepositoryImpl{DB: db}
}

func (p *ProductRepositoryImpl) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
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

func (p *ProductRepositoryImpl) UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
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

func (p *ProductRepositoryImpl) DeleteProduct(ctx context.Context, product *domain.Product) error {
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

func (p *ProductRepositoryImpl) ReadProductById(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	foundProduct := new(domain.Product)
	tx := p.DB.WithContext(ctx)
	if err := tx.Where("product_id = ?", product.ProductID).First(&foundProduct).Error; err != nil {
		return nil, err
	}

	return foundProduct, nil
}

func (p *ProductRepositoryImpl) ReadAllProduct(ctx context.Context, product *domain.Product) ([]domain.Product, error) {
	var products []domain.Product

	tx := p.DB.WithContext(ctx)
	if err := tx.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (p *ProductRepositoryImpl) ReadAllProductSlug(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	foundProduct := new(domain.Product)
	tx := p.DB.WithContext(ctx)
	if err := tx.Where("slug = ?", product.Slug).First(&foundProduct).Error; err != nil {
		return nil, err
	}

	return foundProduct, nil
}
