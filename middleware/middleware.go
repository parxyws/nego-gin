package middleware

import (
	"github.com/parxyws/nego-gin/config"
	"github.com/sirupsen/logrus"
)

type ConfigMiddleware struct {
	Logger *logrus.Logger
	Config *config.Config
}

type ManagerMiddleware struct {
	logger *logrus.Logger
	cfg    *config.Config
}

func NewMiddlewareManager(cfg *ConfigMiddleware) *ManagerMiddleware {
	return &ManagerMiddleware{cfg: cfg.Config, logger: cfg.Logger}
}
