package server

import (
	"github.com/parxyws/nego-gin/internal/middleware"
	productModule "github.com/parxyws/nego-gin/internal/product/module"
	userModule "github.com/parxyws/nego-gin/internal/user/module"
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
	corsMiddleware := middlewareSetup.CORSMiddleware()
	requestIdMiddleware := middlewareSetup.RequestIDMiddleware()

	//------------------ Modules ------------------
	api := s.app.Group("/api/v1")
	api.Use(corsMiddleware)
	api.Use(requestIdMiddleware)

	userModule.InitUserModule(s.db, api, s.cfg, s.rds, s.mail, authMiddleware)
	productModule.InitProductModule(s.db, api, authMiddleware)

	/* ----------------------------- Seed ---------------------------- */
	seeder := NewSeeder(s.db, s.cfg, s.app)
	if err := seeder.Seed(); err != nil {
		return err
	}

	/* ------------------- Permission Cache (startup) ----------------- */
	// Load all role→permission mappings into memory once so that
	// RequirePermission middleware never hits the DB on the hot path.
	if err := middlewareSetup.LoadPermissionCache(); err != nil {
		return err
	}

	return nil
}
