package main

import (
	"log"

	"github.com/parxyws/nego-gin/config"
	"github.com/parxyws/nego-gin/internal/server"
	"github.com/parxyws/nego-gin/pkg/database/aws"
	"github.com/parxyws/nego-gin/pkg/database/psql"
	"github.com/parxyws/nego-gin/pkg/database/rds"
	"github.com/parxyws/nego-gin/pkg/logger"
	"github.com/parxyws/nego-gin/pkg/mail"
)

func main() {
	cfg, err := config.NewAppConfig()
	if err != nil {
		panic(err)
	}

	db, err := psql.NewDB(cfg)
	if err != nil {
		panic(err)
	}

	minio, err := aws.NewAWSClient(cfg)
	if err != nil {
		panic(err)
	}

	authRedis := rds.NewAuthRedis(cfg)
	sessionRedis := rds.NewSessionRedis(cfg)
	limiterRedis := rds.NewLimiterRedis(cfg)

	logrus := logger.NewLogrusLogger(cfg)

	gomail := mail.NewGoMailDialer(cfg)

	s := server.NewServer(&server.ServerConfig{
		Cfg: cfg,
		Db:  db,
		Rds: &rds.RedisClient{
			AuthRedis:    authRedis,
			SessionRedis: sessionRedis,
			LimiterRedis: limiterRedis,
		},
		AwsClient: minio,
		Logger:    logrus,
		Mail:      gomail,
	})

	if err := s.Init(); err != nil {
		log.Fatal(err)
	}

}
