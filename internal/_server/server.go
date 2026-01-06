package _server

import (
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/parxyws/nego-gin/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ServerConfig struct {
	cfg       *config.Config
	db        *gorm.DB
	awsClient *minio.Client
	rds       *redis.Client
}

type Server struct {
	cfg       *config.Config
	app       *gin.Engine
	db        *gorm.DB
	awsClient *minio.Client
	rds       *redis.Client
}

func NewServer(config *ServerConfig) *Server {
	return &Server{
		cfg:       config.cfg,
		db:        config.db,
		awsClient: config.awsClient,
		rds:       config.rds,
	}
}

func (s *Server) Init() error {
	s.app = gin.New()

	return nil
}
