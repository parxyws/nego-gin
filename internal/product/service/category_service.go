package service

import (
	"context"

	"github.com/parxyws/nego-gin/internal/product"
	"github.com/parxyws/nego-gin/internal/product/domain"
	"github.com/parxyws/nego-gin/internal/product/domain/dto"
)

type CategoryServiceImpl struct {
	ctRepo  product.CategoryRepository
	prdRepo product.ProductRepository
}

func NewCategoryServiceImpl(ctRepo product.CategoryRepository, prdRepo product.ProductRepository) product.CategoryService {
	return &CategoryServiceImpl{ctRepo: ctRepo, prdRepo: prdRepo}
}

func (c *CategoryServiceImpl) ListCategories(ctx context.Context) ([]dto.CategoryResponse, error) {
	result, err := c.ctRepo.ReadAllCategory(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]dto.CategoryResponse, len(result))
	for i, category := range result {
		response[i] = dto.CategoryResponse{
			CategoryID:       category.CategoryID,
			ParentCategoryID: category.ParentCategoryID,
			CategoryName:     category.CategoryName,
			Slug:             category.Slug,
			Description:      category.Description,
			ImageURL:         category.ImageURL,
			IsActive:         category.IsActive,
			CreatedAt:        category.CreatedAt,
		}
	}

	return response, nil
}

func (c *CategoryServiceImpl) GetCategory(ctx context.Context, categoryId int32) (*dto.CategoryResponse, error) {
	result, err := c.ctRepo.ReadCategoryById(ctx, &domain.Category{CategoryID: categoryId})
	if err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		CategoryID:       result.CategoryID,
		ParentCategoryID: result.ParentCategoryID,
		CategoryName:     result.CategoryName,
		Slug:             result.Slug,
		Description:      result.Description,
		ImageURL:         result.ImageURL,
		IsActive:         result.IsActive,
		CreatedAt:        result.CreatedAt,
	}, nil
}

func (c *CategoryServiceImpl) ListProductsInCategory(ctx context.Context, categoryId int32) (*dto.CategoryResponse, error) {
	result, err := c.ctRepo.ReadCategoryProductById(ctx, &domain.Category{CategoryID: categoryId})
	if err != nil {
		return nil, err
	}

	products := make([]dto.ProductResponse, len(result.Products))
	for i, p := range result.Products {
		products[i] = dto.ProductResponse{
			ProductID:         p.ProductID,
			CategoryID:        p.CategoryID,
			SellerID:          p.SellerID,
			ProductType:       p.ProductType,
			Sku:               p.Sku,
			Name:              p.Name,
			Slug:              p.Slug,
			Description:       p.Description,
			ShortDescription:  p.ShortDescription,
			BasePrice:         p.BasePrice,
			SalePrice:         p.SalePrice,
			CostPrice:         p.CostPrice,
			StockQuantity:     p.StockQuantity,
			LowStockThreshold: p.LowStockThreshold,
			IsUnlimitedStock:  p.IsUnlimitedStock,
			Status:            p.Status,
			IsFeatured:        p.IsFeatured,
			WeightKg:          p.WeightKg,
			DimensionsCm:      p.DimensionsCm,
			Brand:             p.Brand,
			Condition:         p.Condition,
			ShopCategoryID:    p.ShopCategoryID,
			CreatedAt:         p.CreatedAt,
			UpdatedAt:         p.UpdatedAt,
		}
	}

	return &dto.CategoryResponse{
		CategoryID:       result.CategoryID,
		ParentCategoryID: result.ParentCategoryID,
		CategoryName:     result.CategoryName,
		Slug:             result.Slug,
		Description:      result.Description,
		ImageURL:         result.ImageURL,
		IsActive:         result.IsActive,
		Products:         products,
		CreatedAt:        result.CreatedAt,
	}, nil
}

func (c *CategoryServiceImpl) CreateCategory(ctx context.Context, entity *dto.CategoryCreateRequest) (*dto.CategoryResponse, error) {
	request := &domain.Category{
		ParentCategoryID: entity.ParentCategoryID,
		CategoryName:     entity.CategoryName,
		Description:      entity.Description,
	}

	category, err := c.ctRepo.CreateCategory(ctx, request)
	if err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		CategoryID:       category.CategoryID,
		ParentCategoryID: category.ParentCategoryID,
		CategoryName:     category.CategoryName,
		Slug:             category.Slug,
		Description:      category.Description,
		ImageURL:         category.ImageURL,
		IsActive:         category.IsActive,
		CreatedAt:        category.CreatedAt,
	}, nil
}

func (c *CategoryServiceImpl) UpdateCategory(ctx context.Context, entity *dto.CategoryUpdateRequest) (*dto.CategoryResponse, error) {
	request := &domain.Category{
		CategoryID:       entity.CategoryID,
		ParentCategoryID: entity.ParentCategoryID,
		CategoryName:     entity.CategoryName,
		Description:      entity.Description,
	}

	category, err := c.ctRepo.UpdateCategory(ctx, request)
	if err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		CategoryID:       category.CategoryID,
		ParentCategoryID: category.ParentCategoryID,
		CategoryName:     category.CategoryName,
		Slug:             category.Slug,
		Description:      category.Description,
		ImageURL:         category.ImageURL,
		IsActive:         category.IsActive,
		CreatedAt:        category.CreatedAt,
	}, nil
}

func (c *CategoryServiceImpl) RemoveCategory(ctx context.Context, categoryId int32) error {
	request := &domain.Category{CategoryID: categoryId}

	if err := c.ctRepo.DeleteCategory(ctx, request); err != nil {
		return err
	}

	return nil
}
