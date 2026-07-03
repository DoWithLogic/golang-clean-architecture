package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisPrefixKey string

const (
	REDIS_PREFIX_KEY_CONFIG RedisPrefixKey = "config:%s"
	REDIS_PREFIX_KEY_TOKEN  RedisPrefixKey = "token:%s"
)

const REDIS_TOKEN_EXPIRATION_TIME = time.Minute * 60

func (rpk RedisPrefixKey) String() string { return string(rpk) }

type RedisManager interface {
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (data string, err error)
	Del(ctx context.Context, key string) error
	Close() error
}

// redisManager is a concrete implementation of RedisClient
type redisManager struct {
	client *redis.Client
}

func NewRedisManager(client *redis.Client) RedisManager {
	return &redisManager{
		client: client,
	}
}

// Set sets a value in Redis with a specified expiration.
// Duration to 0 will set expired to 21600 (6 hours)
func (r *redisManager) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	if expiration == 0 {
		expiration = 21600
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
