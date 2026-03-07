package middleware

import (
	"sync"

	"github.com/parxyws/nego-gin/config"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ConfigMiddleware struct {
	Logger *logrus.Logger
	Config *config.Config
	DB     *gorm.DB
	Rdb    *redis.Client
}

type ManagerMiddleware struct {
	logger *logrus.Logger
	cfg    *config.Config
	db     *gorm.DB
	mu     sync.RWMutex
	// permissionCache maps role name → set of permission Names (e.g. "admin:read").
	// Populated once at startup via LoadPermissionCache(); read-only afterwards.
	permissionCache map[string]map[string]bool
}

func NewMiddlewareManager(cfg *ConfigMiddleware) *ManagerMiddleware {
	return &ManagerMiddleware{
		cfg:             cfg.Config,
		logger:          cfg.Logger,
		db:              cfg.DB,
		permissionCache: make(map[string]map[string]bool),
	}
}
