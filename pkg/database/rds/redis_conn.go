package rds

import (
	"fmt"

	"github.com/parxyws/nego-gin/config"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	AuthRedis    *redis.Client
	SessionRedis *redis.Client
	LimiterRedis *redis.Client
}

func NewAuthRedis(config *config.Config) *redis.Client {
	dsn := fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     dsn,                   // Redis server address
		Password: config.Redis.Password, // No password set
		DB:       config.Redis.AuthDb,   // Use default DB
	})

	return rdb
}

func NewSessionRedis(config *config.Config) *redis.Client {
	dsn := fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     dsn,                    // Redis server address
		Password: config.Redis.Password,  // No password set
		DB:       config.Redis.SessionDb, // Use default DB
	})

	return rdb
}

func NewLimiterRedis(config *config.Config) *redis.Client {
	dsn := fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     dsn,                    // Redis server address
		Password: config.Redis.Password,  // No password set
		DB:       config.Redis.LimiterDb, // Use default DB
	})

	return rdb
}
