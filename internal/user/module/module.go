package module

import (
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/config"
	"github.com/parxyws/nego-gin/internal/user/controller"
	"github.com/parxyws/nego-gin/internal/user/repository"
	"github.com/parxyws/nego-gin/internal/user/route"
	"github.com/parxyws/nego-gin/internal/user/service"
	"github.com/redis/go-redis/v9"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

func InitUserModule(
	db *gorm.DB,
	api *gin.RouterGroup,
	cfg *config.Config,
	rds *redis.Client,
	mail *gomail.Dialer,
	authMiddleware *jwt.GinJWTMiddleware,
) {
	// Repositories
	roleRepository := repository.NewRoleRepository(db)
	userRepository := repository.NewUserRepository(db)
	userRoleRepository := repository.NewUserRoleRepository(db)
	userAddressRepository := repository.NewUserAddressRepository(db)

	// Services
	authService := service.NewAuthService(cfg, userRepository, userRoleRepository, roleRepository, rds, mail, authMiddleware)
	userService := service.NewUserService(cfg, userRepository)
	userAddressService := service.NewUserAddressService(userAddressRepository)

	// Controllers
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService, userAddressService)
	adminController := controller.NewAdminController(userService)

	// Routes
	route.AuthRoute(api, authController, authMiddleware)
	route.UserRoute(api, userController, authMiddleware)
	route.AdminRoute(api, adminController, authMiddleware)
}
