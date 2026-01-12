package server

import (
	//------------------ Repository ------------------
	userRepository "github.com/parxyws/nego-gin/internal/user/repository"

	//------------------ Service ------------------
	userService "github.com/parxyws/nego-gin/internal/user/service"

	//------------------ Controller ------------------
	userController "github.com/parxyws/nego-gin/internal/user/controller"

	//------------------ Route ------------------
	userRoute "github.com/parxyws/nego-gin/internal/user/route"

	"github.com/parxyws/nego-gin/middleware"
)

func (s *Server) Boostrap() error {
	middlewareSetup := middleware.NewMiddlewareManager(&middleware.ConfigMiddleware{
		Config: s.cfg,
		Logger: s.logger,
	})

	/* ----------------------------- Repository ---------------------------- */
	RoleRepository := userRepository.NewRoleRepository(s.db)
	UserRepository := userRepository.NewUserRepository(s.db)
	//UserAddressRepository := userRepository.NewUserAddressRepository(s.db)
	UserRoleRepository := userRepository.NewUserRoleRepository(s.db)

	/* ----------------------------- Service ---------------------------- */
	AuthService := userService.NewAuthService(s.cfg, UserRepository, UserRoleRepository, RoleRepository, s.rds, s.mail)

	/* ----------------------------- Controller ---------------------------- */
	AuthController := userController.NewAuthController(AuthService)

	/* ----------------------------- Middleware ---------------------------- */
	_, err := middlewareSetup.JWTAuthMiddleware()
	if err != nil {
		return err
	}

	/* ----------------------------- Route ---------------------------- */
	api := s.app.Group("/api/v1")
	userRoute.AuthRoute(api, AuthController)

	// Note: UserController and AdminController need to be initialized when implemented
	// userRoute.UserRoute(apiV1, UserController)
	// userRoute.AdminRoute(apiV1, AdminController)

	/* ----------------------------- Seed ---------------------------- */
	seeder := NewSeeder(s.db, s.cfg)
	if err := seeder.Seed(); err != nil {
		return err
	}

	return nil
}
