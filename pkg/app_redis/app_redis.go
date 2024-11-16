package app_redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis interface {
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (data string, err error)
	Del(ctx context.Context, key string) error
}

// redisManager is a concrete implementation of RedisClient
type redisManager struct {
	client *redis.Client
}

func NewRedis(client *redis.Client) Redis {
	return &redisManager{
		client: client,
	}
}

// Set sets a value in Redis with a specified expiration.
func (r *redisManager) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	if expiration == 0 {
		expiration = -1
	}

	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value from Redis by key.
func (r *redisManager) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Del deletes a key from Redis.
func (r *redisManager) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Close closes the connection to the Redis server.
func (r *redisManager) Close() error {
	return r.client.Close()
}
