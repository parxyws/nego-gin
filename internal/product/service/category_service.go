package service

import (
	"context"

	"github.com/parxyws/nego-gin/internal/product"
	"github.com/parxyws/nego-gin/internal/product/domain"
	"github.com/parxyws/nego-gin/internal/product/domain/dto"
	"github.com/parxyws/nego-gin/pkg/logger"
)

type CategoryService struct {
	ctRepo  product.CategoryRepository
	prdRepo product.ProductRepository
}

func NewCategoryService(ctRepo product.CategoryRepository, prdRepo product.ProductRepository) product.CategoryService {
	return &CategoryService{ctRepo: ctRepo, prdRepo: prdRepo}
}

func (s *CategoryService) ListCategories(ctx context.Context) ([]dto.CategoryResponse, error) {
	log := logger.WithCtx(ctx, "service", "CategoryService.ListCategories")

	categories, err := s.ctRepo.ReadAllCategory(ctx)
	if err != nil {
		log.Errorf("failed to read all categories: %v", err)
		return nil, err
	}

	categoriesResp := make([]dto.CategoryResponse, len(categories))
	for i, category := range categories {
		categoriesResp[i] = dto.CategoryResponse{
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

	log.Info("category listed successfully")
	return categoriesResp, nil
}

func (s *CategoryService) GetCategory(ctx context.Context, categoryID int32) (*dto.CategoryResponse, error) {
	log := logger.WithCtx(ctx, "service", "CategoryService.GetCategory")

	category, err := s.ctRepo.ReadCategoryById(ctx, &domain.Category{CategoryID: categoryID})
	if err != nil {
		log.Errorf("failed to read category by id: %v", err)
		return nil, err
	}

	log.Info("category get successfully")
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

func (s *CategoryService) ListProductsInCategory(ctx context.Context, categoryID int32) (*dto.CategoryResponse, error) {
	log := logger.WithCtx(ctx, "service", "CategoryService.ListProductsInCategory")

	categoryProducts, err := s.ctRepo.ReadCategoryProductById(ctx, &domain.Category{CategoryID: categoryID})
	if err != nil {
		log.Errorf("failed to read category products: %v", err)
		return nil, err
	}

	products := make([]dto.ProductResponse, len(categoryProducts.Products))
	for i, product := range categoryProducts.Products {
		products[i] = dto.ProductResponse{
			ProductID:         product.ProductID,
			CategoryID:        product.CategoryID,
			SellerID:          product.SellerID,
			ProductType:       product.ProductType,
			Sku:               product.Sku,
			Name:              product.Name,
			Slug:              product.Slug,
			Description:       product.Description,
			ShortDescription:  product.ShortDescription,
			BasePrice:         product.BasePrice,
			SalePrice:         product.SalePrice,
			CostPrice:         product.CostPrice,
			StockQuantity:     product.StockQuantity,
			LowStockThreshold: product.LowStockThreshold,
			IsUnlimitedStock:  product.IsUnlimitedStock,
			Status:            product.Status,
			IsFeatured:        product.IsFeatured,
			WeightKg:          product.WeightKg,
			DimensionsCm:      product.DimensionsCm,
			Brand:             product.Brand,
			Condition:         product.Condition,
			ShopCategoryID:    product.ShopCategoryID,
			CreatedAt:         product.CreatedAt,
			UpdatedAt:         product.UpdatedAt,
		}
	}

	log.Info("category get listed successfully")
	return &dto.CategoryResponse{
		CategoryID:       categoryProducts.CategoryID,
		ParentCategoryID: categoryProducts.ParentCategoryID,
		CategoryName:     categoryProducts.CategoryName,
		Slug:             categoryProducts.Slug,
		Description:      categoryProducts.Description,
		ImageURL:         categoryProducts.ImageURL,
		IsActive:         categoryProducts.IsActive,
		Products:         products,
		CreatedAt:        categoryProducts.CreatedAt,
	}, nil
}

func (s *CategoryService) CreateCategory(ctx context.Context, req *dto.CategoryCreateRequest) (*dto.CategoryResponse, error) {
	log := logger.WithCtx(ctx, "service", "CategoryService.CreateCategory").WithField("category_name", req.CategoryName)

	categoryRequest := &domain.Category{
		ParentCategoryID: req.ParentCategoryID,
		CategoryName:     req.CategoryName,
		Description:      req.Description,
	}

	category, err := s.ctRepo.CreateCategory(ctx, categoryRequest)
	if err != nil {
		log.Errorf("failed to create category in repo: %v", err)
		return nil, err
	}

	log.WithField("category_id", category.CategoryID).Info("category created successfully")

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

func (s *CategoryService) UpdateCategory(ctx context.Context, req *dto.CategoryUpdateRequest) (*dto.CategoryResponse, error) {
	log := logger.WithCtx(ctx, "service", "CategoryService.UpdateCategory").WithField("category_id", req.CategoryID)

	categoryRequest := &domain.Category{
		CategoryID:       req.CategoryID,
		ParentCategoryID: req.ParentCategoryID,
		CategoryName:     req.CategoryName,
		Description:      req.Description,
	}

	category, err := s.ctRepo.UpdateCategory(ctx, categoryRequest)
	if err != nil {
		log.Errorf("failed to update category in repo: %v", err)
		return nil, err
	}

	log.Info("category updated successfully")

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

func (s *CategoryService) RemoveCategory(ctx context.Context, categoryID int32) error {
	log := logger.WithCtx(ctx, "service", "CategoryService.RemoveCategory").WithField("category_id", categoryID)
	categoryRequest := &domain.Category{CategoryID: categoryID}

	if err := s.ctRepo.DeleteCategory(ctx, categoryRequest); err != nil {
		log.Errorf("failed to delete category in repo: %v", err)
		return err
	}

	log.Info("category deleted successfully")
	return nil
}
