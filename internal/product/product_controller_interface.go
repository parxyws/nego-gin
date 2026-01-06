package product

import "github.com/gin-gonic/gin"

type CategoryController interface {
	ListCategories(c *gin.Context)
	GetCategory(c *gin.Context)
	ListProductsByCategory(c *gin.Context)
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
}
