package server

import (
	"github.com/parxyws/nego-gin/internal/middleware"
	//------------------ Repository ------------------
	userRepository "github.com/parxyws/nego-gin/internal/user/repository"

	//------------------ Service ------------------
	userService "github.com/parxyws/nego-gin/internal/user/service"

	//------------------ Controller ------------------
	userController "github.com/parxyws/nego-gin/internal/user/controller"

	//------------------ Route ------------------
	userRoute "github.com/parxyws/nego-gin/internal/user/route"
)

func (s *Server) Boostrap() error {
	middlewareSetup := middleware.NewMiddlewareManager(&middleware.ConfigMiddleware{
		Config: s.cfg,
		Logger: s.logger,
		DB:     s.db,
	})

	/* ----------------------------- Middleware ---------------------------- */
	authMiddleware, err := middlewareSetup.JWTAuthMiddleware()
	if err != nil {
		return err
	}

	/* ----------------------------- Repository ---------------------------- */
	//PermissionRepository := userRepository.NewPermissionRepository(s.db)
	RoleRepository := userRepository.NewRoleRepository(s.db)
	//RolePermissionRepository := userRepository.NewRolePermissionRepository(s.db)
	UserRepository := userRepository.NewUserRepository(s.db)
	UserRoleRepository := userRepository.NewUserRoleRepository(s.db)
	UserAddressRepository := userRepository.NewUserAddressRepository(s.db)

	/* ----------------------------- Service ---------------------------- */
	//RoleService := userService.NewRoleService(RoleRepository, UserRoleRepository)
	AuthService := userService.NewAuthService(s.cfg, UserRepository, UserRoleRepository, RoleRepository, s.rds, s.mail, authMiddleware)
	UserService := userService.NewUserService(s.cfg, UserRepository)
	UserAddressService := userService.NewUserAddressService(UserAddressRepository)
	/* ----------------------------- Controller ---------------------------- */
	AuthController := userController.NewAuthController(AuthService)
	UserController := userController.NewUserController(UserService, UserAddressService)

	/* ----------------------------- Route ---------------------------- */
	api := s.app.Group("/api/v1")
	userRoute.AuthRoute(api, AuthController)
	userRoute.UserRoute(api, UserController)

	/* ----------------------------- Seed ---------------------------- */
	seeder := NewSeeder(s.db, s.cfg, s.app)
	if err := seeder.Seed(); err != nil {
		return err
	}

	return nil
}
