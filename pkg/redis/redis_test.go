package redis_test

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"

	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/datasources"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/redis"
)

func TestRedisManager(t *testing.T) {
	// Start an in-memory Redis server
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	// Create a go-redis client connected to miniredis
	client := datasources.NewRedisClient(context.Background(), config.RedisConfig{Addr: mr.Addr(), DB: 0})

	redisManager := redis.NewRedis(client) // renamed to avoid name clash
	ctx := context.Background()

	t.Run("Set and Get", func(t *testing.T) {
		key := "first_key"
		value := "first_value"

		// Test Set
		err := redisManager.Set(ctx, key, value, 0)
		assert.NoError(t, err)

		// Verify stored value in miniredis
		mr.CheckGet(t, key, value)

		// Test Get
		retrieved, err := redisManager.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, retrieved)
	})

	t.Run("Delete key", func(t *testing.T) {
		key := "second_key"
		value := "second_value"

		err := redisManager.Set(ctx, key, value, time.Minute)
		assert.NoError(t, err)

		// Delete key
		err = redisManager.Del(ctx, key)
		assert.NoError(t, err)

		// Ensure its gone
		_, err = redisManager.Get(ctx, key)
		assert.Error(t, err, "Expected error when getting deleted key")
	})

	t.Run("Close client", func(t *testing.T) {
		err := redisManager.Close()
		assert.NoError(t, err)
	})
}
