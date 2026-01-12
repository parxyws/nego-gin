package repository

import (
	"context"

	"github.com/parxyws/nego-gin/internal/product"
	"github.com/parxyws/nego-gin/internal/product/domain"
	"gorm.io/gorm"
)

type ProductMediaRepositoryImpl struct {
	DB *gorm.DB
}

func NewProductMediaRepositoryImpl(db *gorm.DB) product.ProductMediaRepository {
	return &ProductMediaRepositoryImpl{DB: db}
}

func (p *ProductMediaRepositoryImpl) CreateProductMedia(ctx context.Context, product *domain.ProductMedia) (*domain.ProductMedia, error) {
	tx := p.DB.WithContext(ctx)

	err := tx.Transaction(func(tx *gorm.DB) error {
		if result := tx.Create(&product); result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *ProductMediaRepositoryImpl) UpdateProductMedia(ctx context.Context, product *domain.ProductMedia) (*domain.ProductMedia, error) {
	tx := p.DB.WithContext(ctx)
	err := tx.Transaction(func(tx *gorm.DB) error {
		if result := tx.Save(&product); result.Error != nil {
			return result.Error
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *ProductMediaRepositoryImpl) DeleteProductMedia(ctx context.Context, product *domain.ProductMedia) error {
	tx := p.DB.WithContext(ctx)
	err := tx.Transaction(func(tx *gorm.DB) error {
		if result := tx.Delete(&product); result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (p *ProductMediaRepositoryImpl) ReadProductMediaById(ctx context.Context, product *domain.ProductMedia) (*domain.ProductMedia, error) {
	var productMedia *domain.ProductMedia
	tx := p.DB.WithContext(ctx)

	if err := tx.Where("media_id = ?", product.MediaID).First(&productMedia).Error; err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return productMedia, nil
}

func (p *ProductMediaRepositoryImpl) ReadAllProductMediaById(ctx context.Context, product *domain.ProductMedia) ([]domain.ProductMedia, error) {
	var productMediaList []domain.ProductMedia
	tx := p.DB.WithContext(ctx)

	if err := tx.Find(&productMediaList).Error; err != nil {
		return nil, err
	}

	return productMediaList, nil
}
