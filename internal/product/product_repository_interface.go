package product

import (
	"context"

	"github.com/parxyws/nego-gin/internal/product/domain"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error)
	UpdateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error)
	DeleteCategory(ctx context.Context, category *domain.Category) error

	ReadCategoryById(ctx context.Context, category *domain.Category) (*domain.Category, error)
	ReadCategoryProductById(ctx context.Context, category *domain.Category) (*domain.Category, error)
	ReadCategoryByName(ctx context.Context, category *domain.Category) (*domain.Category, error)
	ReadAllCategory(ctx context.Context) ([]domain.Category, error)
}

type ShopCategoryRepository interface {
	CreateShopCategory(ctx context.Context, category *domain.ShopCategory) (*domain.ShopCategory, error)
	UpdateShopCategory(ctx context.Context, category *domain.ShopCategory) (*domain.ShopCategory, error)
	DeleteShopCategory(ctx context.Context, category *domain.ShopCategory) error
	ReadShopCategoryById(ctx context.Context, category *domain.ShopCategory) (*domain.ShopCategory, error)
	ReadShopCategoryProductById(ctx context.Context, category *domain.ShopCategory) (*domain.ShopCategory, error)
	ReadAllShopCategory(ctx context.Context) ([]domain.ShopCategory, error)
}

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	DeleteProduct(ctx context.Context, product *domain.Product) error

	ReadProductById(ctx context.Context, product *domain.Product) (*domain.Product, error)
	ReadAllProduct(ctx context.Context, product *domain.Product) ([]domain.Product, error)
	ReadAllProductSlug(ctx context.Context, product *domain.Product) (*domain.Product, error)
}

type ProductVariantRepository interface {
	CreateProductVariant(ctx context.Context, product *domain.ProductVariant) (*domain.ProductVariant, error)
	UpdateProductVariant(ctx context.Context, product *domain.ProductVariant) (*domain.ProductVariant, error)
	DeleteProductVariant(ctx context.Context, product *domain.ProductVariant) error

	ReadProductVariantById(ctx context.Context, product *domain.ProductVariant) (*domain.ProductVariant, error)
	ReadProductVariantByName(ctx context.Context, product *domain.ProductVariant) (*domain.ProductVariant, error)
	ReadAllProductVariantByProductId(ctx context.Context, product *domain.ProductVariant) ([]domain.ProductVariant, error)
}

type ProductAttributeRepository interface {
	CreateProductAttribute(ctx context.Context, product *domain.ProductAttribute) (*domain.ProductAttribute, error)
	UpdateProductAttribute(ctx context.Context, product *domain.ProductAttribute) (*domain.ProductAttribute, error)
	DeleteProductAttribute(ctx context.Context, product *domain.ProductAttribute) error

	ReadProductAttributeById(ctx context.Context, product *domain.ProductAttribute) (*domain.ProductAttribute, error)
	ReadProductAttributeByName(ctx context.Context, product *domain.ProductAttribute) (*domain.ProductAttribute, error)
	ReadAllProductAttributeByProductId(ctx context.Context, product *domain.ProductAttribute) ([]domain.ProductAttribute, error)
}

type ProductMediaRepository interface {
	CreateProductMedia(ctx context.Context, product *domain.ProductMedia) (*domain.ProductMedia, error)
	UpdateProductMedia(ctx context.Context, product *domain.ProductMedia) (*domain.ProductMedia, error)
	DeleteProductMedia(ctx context.Context, product *domain.ProductMedia) error

	ReadProductMediaById(ctx context.Context, product *domain.ProductMedia) (*domain.ProductMedia, error)
	ReadAllProductMediaById(ctx context.Context, product *domain.ProductMedia) ([]domain.ProductMedia, error)
}

type TagRepository interface {
	CreateTag(ctx context.Context, tag *domain.Tag) (*domain.Tag, error)
	UpdateTag(ctx context.Context, tag *domain.Tag) (*domain.Tag, error)
	DeleteTag(ctx context.Context, tag *domain.Tag) error
	ReadTagById(ctx context.Context, tag *domain.Tag) (*domain.Tag, error)
	ReadAllTag(ctx context.Context) ([]domain.Tag, error)
}

type ProductTagRepository interface {
	CreateProductTag(ctx context.Context, product *domain.ProductTag) (*domain.ProductTag, error)
	UpdateProductTag(ctx context.Context, product *domain.ProductTag) (*domain.ProductTag, error)
	DeleteProductTag(ctx context.Context, product *domain.ProductTag) error

	ReadProductTagByTagId(ctx context.Context, product *domain.ProductTag) (*domain.ProductTag, error)
	ReadAllProductTagByProductId(ctx context.Context, product *domain.ProductTag) ([]domain.ProductTag, error)
}
