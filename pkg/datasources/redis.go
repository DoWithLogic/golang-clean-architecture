package datasources

import (
	"context"

	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/go-redis/redis/v8"
)

func NewRedisClient(ctx context.Context, cfg config.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		panic(err)
	}

	return client
}
