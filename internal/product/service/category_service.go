package service

import (
	"context"

	"github.com/parxyws/nego-gin/internal/product"
	"github.com/parxyws/nego-gin/internal/product/domain"
	"github.com/parxyws/nego-gin/internal/product/domain/dto"
	"github.com/sirupsen/logrus"
)

type CategoryServiceImpl struct {
	ctRepo  product.CategoryRepository
	prdRepo product.ProductRepository
}

func NewCategoryServiceImpl(ctRepo product.CategoryRepository, prdRepo product.ProductRepository) product.CategoryService {
	return &CategoryServiceImpl{ctRepo: ctRepo, prdRepo: prdRepo}
}

func (s *CategoryServiceImpl) ListCategories(ctx context.Context) ([]dto.CategoryResponse, error) {
	categories, err := s.ctRepo.ReadAllCategory(ctx)
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "CategoryService.ListCategories"}).Errorf("failed to read all categories: %v", err)
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

	return categoriesResp, nil
}

func (s *CategoryServiceImpl) GetCategory(ctx context.Context, categoryID int32) (*dto.CategoryResponse, error) {
	category, err := s.ctRepo.ReadCategoryById(ctx, &domain.Category{CategoryID: categoryID})
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "CategoryService.GetCategory", "category_id": categoryID}).Errorf("failed to read category by id: %v", err)
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

func (s *CategoryServiceImpl) ListProductsInCategory(ctx context.Context, categoryID int32) (*dto.CategoryResponse, error) {
	categoryProducts, err := s.ctRepo.ReadCategoryProductById(ctx, &domain.Category{CategoryID: categoryID})
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "CategoryService.ListProductsInCategory", "category_id": categoryID}).Errorf("failed to read category products: %v", err)
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

func (s *CategoryServiceImpl) CreateCategory(ctx context.Context, req *dto.CategoryCreateRequest) (*dto.CategoryResponse, error) {
	categoryRequest := &domain.Category{
		ParentCategoryID: req.ParentCategoryID,
		CategoryName:     req.CategoryName,
		Description:      req.Description,
	}

	category, err := s.ctRepo.CreateCategory(ctx, categoryRequest)
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "CategoryService.CreateCategory", "category_name": req.CategoryName}).Errorf("failed to create category in repo: %v", err)
		return nil, err
	}

	logrus.WithFields(logrus.Fields{"function": "CategoryService.CreateCategory", "category_id": category.CategoryID}).Info("category created successfully")

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

func (s *CategoryServiceImpl) UpdateCategory(ctx context.Context, req *dto.CategoryUpdateRequest) (*dto.CategoryResponse, error) {
	categoryRequest := &domain.Category{
		CategoryID:       req.CategoryID,
		ParentCategoryID: req.ParentCategoryID,
		CategoryName:     req.CategoryName,
		Description:      req.Description,
	}

	category, err := s.ctRepo.UpdateCategory(ctx, categoryRequest)
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "CategoryService.UpdateCategory", "category_id": req.CategoryID}).Errorf("failed to update category in repo: %v", err)
		return nil, err
	}

	logrus.WithFields(logrus.Fields{"function": "CategoryService.UpdateCategory", "category_id": category.CategoryID}).Info("category updated successfully")

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

func (s *CategoryServiceImpl) RemoveCategory(ctx context.Context, categoryID int32) error {
	categoryRequest := &domain.Category{CategoryID: categoryID}

	if err := s.ctRepo.DeleteCategory(ctx, categoryRequest); err != nil {
		logrus.WithFields(logrus.Fields{"function": "CategoryService.RemoveCategory", "category_id": categoryID}).Errorf("failed to delete category in repo: %v", err)
		return err
	}

	logrus.WithFields(logrus.Fields{"function": "CategoryService.RemoveCategory", "category_id": categoryID}).Info("category deleted successfully")
	return nil
}
