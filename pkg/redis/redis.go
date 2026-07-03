package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Addr     string // The address of the database.
	Password string // The password for connecting to the database.
	DB       int    // The name of the database.
}

func NewRedisClient(ctx context.Context, cfg RedisConfig) *redis.Client {
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
