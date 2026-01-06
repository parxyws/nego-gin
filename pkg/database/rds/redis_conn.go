package rds

import (
	"fmt"

	"github.com/parxyws/nego-gin/config"
	"github.com/redis/go-redis/v9"
)

func NewRedis(config *config.Config) *redis.Client {
	dsn := fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     dsn,                   // Redis server address
		Password: config.Redis.Password, // No password set
		DB:       config.Redis.Db,       // Use default DB
	})

	return rdb
}
