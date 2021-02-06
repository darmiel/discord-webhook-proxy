package dbredis

import (
	"github.com/darmiel/whgoxy/internal/whgoxy/config"
	"github.com/go-redis/redis/v8"
)

func NewClient(cfg config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}

var GlobalRedis *redis.Client
