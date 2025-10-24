package datasources_test

import (
	"context"
	"testing"

	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/datasources"
	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
)

func TestNewRedisClient(t *testing.T) {
	// Start a miniredis instance
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer mr.Close()

	ctx := context.Background()

	// Mock Redis configuration
	cfg := config.RedisConfig{
		Addr:     mr.Addr(), // Use the miniredis address
		Password: "",        // No password for miniredis
		DB:       0,         // Default database
	}

	// Call the NewRedisClient function
	client := datasources.NewRedisClient(ctx, cfg)
	assert.NotNil(t, client, "Expected non-nil Redis client")

	// Perform a sample Redis operation to verify the client works
	err = client.Set(ctx, "test_key", "test_value", 0).Err()
	assert.NoError(t, err, "Expected no error setting a key in Redis")

	value, err := client.Get(ctx, "test_key").Result()
	assert.NoError(t, err, "Expected no error getting a key from Redis")
	assert.Equal(t, "test_value", value, "Expected value to match the set value")

	// Simulate a Redis server down scenario
	mr.Close()

	// Try to ping the closed Redis server
	err = client.Ping(ctx).Err()
	assert.Error(t, err, "Expected an error when pinging a closed Redis server")
}
