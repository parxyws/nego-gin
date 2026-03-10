package repository

import (
	"context"

	"github.com/parxyws/nego-gin/internal/product"
	"github.com/parxyws/nego-gin/internal/product/domain"
	"gorm.io/gorm"
)

type ProductMediaRepository struct {
	DB *gorm.DB
}

func NewProductMediaRepository(db *gorm.DB) product.ProductMediaRepository {
	return &ProductMediaRepository{DB: db}
}

func (p *ProductMediaRepository) CreateProductMedia(ctx context.Context, product *domain.ProductMedia) (*domain.ProductMedia, error) {
	tx := p.DB.WithContext(ctx)

	if result := tx.Create(&product); result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

func (p *ProductMediaRepository) UpdateProductMedia(ctx context.Context, product *domain.ProductMedia) (*domain.ProductMedia, error) {
	tx := p.DB.WithContext(ctx)
	if result := tx.Model(&domain.ProductMedia{}).Where("media_id = ?", product.MediaID).Updates(&product); result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

func (p *ProductMediaRepository) DeleteProductMedia(ctx context.Context, product *domain.ProductMedia) error {
	tx := p.DB.WithContext(ctx)
	if result := tx.Delete(&product); result.Error != nil {
		return result.Error
	}

	return nil
}

func (p *ProductMediaRepository) ReadProductMediaById(ctx context.Context, product *domain.ProductMedia) (*domain.ProductMedia, error) {
	var productMedia *domain.ProductMedia
	tx := p.DB.WithContext(ctx)

	if err := tx.Where("media_id = ?", product.MediaID).First(&productMedia).Error; err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return productMedia, nil
}

func (p *ProductMediaRepository) ReadAllProductMediaById(ctx context.Context, product *domain.ProductMedia) ([]domain.ProductMedia, error) {
	var productMediaList []domain.ProductMedia
	tx := p.DB.WithContext(ctx)

	if err := tx.Find(&productMediaList).Error; err != nil {
		return nil, err
	}

	return productMediaList, nil
}
