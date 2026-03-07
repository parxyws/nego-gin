package module

import (
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/product/controller"
	"github.com/parxyws/nego-gin/internal/product/repository"
	"github.com/parxyws/nego-gin/internal/product/route"
	"github.com/parxyws/nego-gin/internal/product/service"
	"gorm.io/gorm"
)

func InitProductModule(db *gorm.DB, api *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	// Repositories
	categoryRepository := repository.NewCategoryRepository(db)
	productRepository := repository.NewProductRepository(db)

	// Services
	categoryService := service.NewCategoryService(categoryRepository, productRepository)

	// Controllers
	categoryController := controller.NewCategoryController(categoryService)

	// Routes — pass authMiddleware for protected endpoints
	route.CategoryRoute(api, categoryController, authMiddleware)
}
