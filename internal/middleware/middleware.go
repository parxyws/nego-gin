package middleware

import (
	"github.com/parxyws/nego-gin/config"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ConfigMiddleware struct {
	Logger *logrus.Logger
	Config *config.Config
	DB     *gorm.DB
}

type ManagerMiddleware struct {
	logger *logrus.Logger
	cfg    *config.Config
	db     *gorm.DB
}

func NewMiddlewareManager(cfg *ConfigMiddleware) *ManagerMiddleware {
	return &ManagerMiddleware{cfg: cfg.Config, logger: cfg.Logger, db: cfg.DB}
}
