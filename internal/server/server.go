package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/parxyws/nego-gin/config"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

const (
	ctxTimeout = 5
	certFile   = "./certs/server.crt"
	keyFile    = "./certs/server.key"
)

type ServerConfig struct {
	Cfg       *config.Config
	Db        *gorm.DB
	AwsClient *minio.Client
	Rds       *redis.Client
	Logger    *logrus.Logger
	Mail      *gomail.Dialer
}

type Server struct {
	cfg       *config.Config
	app       *gin.Engine
	db        *gorm.DB
	awsClient *minio.Client
	rds       *redis.Client
	logger    *logrus.Logger
	mail      *gomail.Dialer
}

func NewServer(config *ServerConfig) *Server {
	return &Server{
		cfg:       config.Cfg,
		db:        config.Db,
		awsClient: config.AwsClient,
		rds:       config.Rds,
		logger:    config.Logger,
		mail:      config.Mail,
	}
}

func (s *Server) Init() error {

	if s.cfg.Server.Mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	} else if s.cfg.Server.Mode == gin.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.TestMode)
	}

	s.app = gin.New()

	if err := s.Boostrap(); err != nil {
		return err
	}

	addr := fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.Port)

	srv := &http.Server{
		Addr:    addr,
		Handler: s.app,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if s.cfg.Server.SSL {
		serverError := make(chan error)

		go func() {
			serverError <- srv.ListenAndServeTLS(certFile, keyFile)
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		select {
		case err := <-serverError:
			log.Fatalf("server error: %s", err)
		case <-quit:
			if err := srv.Shutdown(ctx); err != nil {
				log.Fatalf("server shutdown error: %s", err)
			}
			log.Println("server shutdown gracefully")
			return nil
		}

	} else {
		serverError := make(chan error)

		go func() {
			serverError <- srv.ListenAndServe()
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

		select {
		case err := <-serverError:
			log.Fatalf("Failed to start TLS server: %v", err)
		case <-quit:
			if err := srv.Shutdown(ctx); err != nil {
				log.Fatalf("Error gracefully shutting down server: %v", err)
			}
			log.Println("Server exited properly")
			return nil
		}
	}

	return nil
}
