package redis_test

import (
	"context"
	"testing"

	"github.com/DoWithLogic/golang-clean-architecture/pkg/redis"
	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
)

func TestRedisManager(t *testing.T) {
	// Start a miniredis server
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer mr.Close()

	// Create a Redis client connected to miniredis
	redisManager := redis.NewRedisManager(redis.NewRedisClient(context.Background(), redis.RedisConfig{Addr: mr.Addr()}))

	ctx := context.Background()

	t.Run("Set and Get", func(t *testing.T) {
		key := "test_key"
		value := "test_value"

		// Test setting a value
		err := redisManager.Set(ctx, key, value, 0) // Default expiration
		assert.NoError(t, err)

		// Verify the value is stored in Redis
		mr.CheckGet(t, key, value)

		// Test getting the value
		retrievedValue, err := redisManager.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, retrievedValue)
	})

	t.Run("Delete key", func(t *testing.T) {
		key := "deletable_key"
		value := "to_be_deleted"

		// Set a value to delete
		err := redisManager.Set(ctx, key, value, 0)
		assert.NoError(t, err)

		// Delete the key
		err = redisManager.Del(ctx, key)
		assert.NoError(t, err)

		// Verify the key no longer exists
		_, err = redisManager.Get(ctx, key)
		assert.Error(t, err)
	})

	t.Run("Close client", func(t *testing.T) {
		// Test closing the Redis client
		err := redisManager.Close()
		assert.NoError(t, err)
	})
}
