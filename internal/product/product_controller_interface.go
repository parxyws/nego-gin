package product

import "github.com/gin-gonic/gin"

type CategoryController interface {
	ListCategories(c *gin.Context)
	GetCategory(c *gin.Context)
	ListProductsInCategory(c *gin.Context)
	CreateCategory(c *gin.Context)
	UpdateCategory(c *gin.Context)
	RemoveCategory(c *gin.Context)
}

type TagController interface {
	ListTags(c *gin.Context)
	GetProductsByTag(c *gin.Context)
}

type ProductController interface {
	ListProducts(c *gin.Context)
	GetProductDetails(c *gin.Context)
	GetProductSlug(c *gin.Context)
	CreateProduct(c *gin.Context)
	UpdateProduct(c *gin.Context)
	DeleteProduct(c *gin.Context)
	UpdateProductStatus(c *gin.Context)
	ListSellerProducts(c *gin.Context)
}

type ProductVariantsController interface {
	ListProductVariants(c *gin.Context)
	CreateProductVariant(c *gin.Context)
	UpdateProductVariant(c *gin.Context)
	DeleteProductVariant(c *gin.Context)
}

type ProductMediaController interface {
	ListProductMedia(c *gin.Context)
	UploadMedia(c *gin.Context)
	UpdateProductMedia(c *gin.Context)
	DeleteProductMedia(c *gin.Context)
}
