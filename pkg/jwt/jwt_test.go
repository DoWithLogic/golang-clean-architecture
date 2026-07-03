package jwt_test

import (
	"context"
	"testing"
	"time"

	pkgJWT "github.com/DoWithLogic/golang-clean-architecture/pkg/jwt"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/redis"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/response"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/response/app_error"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/types"

	"github.com/alicebob/miniredis"
	"github.com/go-faker/faker/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupJWTFactory(t *testing.T) (*pkgJWT.JWTFactory, pkgJWT.JWTConfig, func()) {
	t.Helper()

	// Start in-memory Redis
	mr, err := miniredis.Run()
	require.NoError(t, err)

	cfg := pkgJWT.JWTConfig{
		Key:             faker.UUIDDigit(),
		ExpiredInSecond: 3600,
	}

	// Initialize Redis client
	redisClient := redis.NewRedisClient(context.Background(), redis.RedisConfig{Addr: mr.Addr()})
	securityFactory := pkgJWT.NewJWTFactory(cfg, redis.NewRedisManager(redisClient))

	// Cleanup function
	cleanup := func() {
		mr.Close()
	}

	return securityFactory, cfg, cleanup
}

func generateTestClaims(jwtConfig *pkgJWT.JWTConfig) *pkgJWT.JWTClaims {
	expiredTime := time.Now().Add(time.Minute * time.Duration(jwtConfig.ExpiredInSecond))
	return &pkgJWT.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredTime),
		},
		Data: &pkgJWT.Data{
			ID:           1,
			ContactType:  types.CONTACT_TYPE_EMAIL,
			ContactValue: "golangcleanarchitecture@yopmail.com",
		},
	}
}

func TestJWTFactory(t *testing.T) {
	securityFactory, jwtConfig, cleanup := setupJWTFactory(t)
	defer cleanup()

	testCases := []struct {
		name      string
		setup     func() (string, *pkgJWT.JWTClaims, error)
		expectErr bool
	}{
		{
			name: "success verify jwt",
			setup: func() (string, *pkgJWT.JWTClaims, error) {
				requestClaims := generateTestClaims(&jwtConfig)
				bearerToken, err := securityFactory.CreateJWT(requestClaims)
				return bearerToken, requestClaims, err
			},
			expectErr: false,
		},
		{
			name: "request using revoked token should be err",
			setup: func() (string, *pkgJWT.JWTClaims, error) {
				requestClaims := generateTestClaims(&jwtConfig)
				bearerToken, err := securityFactory.CreateJWT(requestClaims)
				if err != nil {
					return "", nil, err
				}
				if err := securityFactory.AddToBlacklist(context.TODO(), bearerToken, requestClaims.ExpiresAt.Time); err != nil {
					return "", nil, err
				}
				return bearerToken, requestClaims, nil
			},
			expectErr: true,
		},
		{
			name: "parse with claims error",
			setup: func() (string, *pkgJWT.JWTClaims, error) {
				return faker.UUIDDigit(), nil, nil
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bearerToken, expectedClaims, err := tc.setup()
			assert.NoError(t, err)
			assert.NotEmpty(t, bearerToken)

			claims, err := securityFactory.VerifyJWT(context.Background(), bearerToken)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Empty(t, claims)
				assert.Equal(t, err, response.Unauthorized(app_error.ErrInvalidToken))
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, claims)
				assert.Equal(t, expectedClaims, claims)
			}
		})
	}
}
