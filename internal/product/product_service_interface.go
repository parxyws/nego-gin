package product

import (
	"context"

	"github.com/parxyws/nego-gin/internal/product/domain/dto"
)

type CategoryService interface {
	ListCategories(ctx context.Context) ([]dto.CategoryResponse, error)
	GetCategory(ctx context.Context, categoryId int32) (*dto.CategoryResponse, error)
	ListProductsInCategory(ctx context.Context, categoryId int32) (*dto.CategoryResponse, error)
	CreateCategory(ctx context.Context, entity *dto.CategoryCreateRequest) (*dto.CategoryResponse, error)
	UpdateCategory(ctx context.Context, entity *dto.CategoryUpdateRequest) (*dto.CategoryResponse, error)
	RemoveCategory(ctx context.Context, categoryId int32) error
}

type ShopCategoryService interface {
	ListShopCategories(ctx context.Context, sellerId string) ([]dto.ShopCategoryResponse, error)
	GetShopCategory(ctx context.Context, categoryId int32) (*dto.ShopCategoryResponse, error)
	CreateShopCategory(ctx context.Context, entity *dto.ShopCategoryCreateRequest) (*dto.ShopCategoryResponse, error)
	UpdateShopCategory(ctx context.Context, entity *dto.ShopCategoryUpdateRequest) (*dto.ShopCategoryResponse, error)
	RemoveShopCategory(ctx context.Context, categoryId int32) error
}

type TagService interface {
	ListTags(ctx context.Context, categoryId int32) ([]dto.TagResponse, error)
	GetProductsByTag(ctx context.Context, tagId int32) ([]dto.ProductResponse, error)
}

type ProductService interface {
	ListProducts(ctx context.Context, sellerId string) ([]dto.ProductResponse, error)
	GetProductDetails(ctx context.Context, productId string) (*dto.ProductResponse, error)
	GetProductSlug(ctx context.Context, slug string) (*dto.ProductResponse, error)
	CreateProduct(ctx context.Context, entity *dto.ProductCreateRequest) (*dto.ProductResponse, error)
	UpdateProduct(ctx context.Context, entity *dto.ProductUpdateRequest) (*dto.ProductResponse, error)
	DeleteProduct(ctx context.Context, productId string) error
	UpdateProductStatus(ctx context.Context, productId string, status string) (*dto.ProductResponse, error)
	ListSellerProducts(ctx context.Context, sellerId string) ([]dto.ProductResponse, error)
}
