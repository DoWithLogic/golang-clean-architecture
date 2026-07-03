package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/DoWithLogic/golang-clean-architecture/pkg/redis"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/response"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/response/app_error"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/types"
	"github.com/golang-jwt/jwt/v5"
)

// Data struct holds the user-related information that is embedded in the JWT token claims.
type Data struct {
	ID           int64              `json:"id"`
	ContactType  types.CONTACT_TYPE `json:"contact_type"`
	ContactValue string             `json:"contact_value"`
}

type JWTConfig struct {
	Key             string
	ExpiredInSecond int64
}

// JWTClaims defines the structure of the data stored in the JWT token.
// It includes standard registered claims along with the custom data field for user-specific information.
type JWTClaims struct {
	jwt.RegisteredClaims       // Includes standard claims like Issuer, Audience, ExpirationTime, etc.
	Data                 *Data `json:"data"` // Custom data field holding the user's information
}

// JWTFactory struct contains the configuration for JWT token creation and verification.
// It interacts with Redis to handle blacklisting of JWT tokens.
type JWTFactory struct {
	cfg   JWTConfig          // Configuration for JWT token generation (e.g., secret key)
	redis redis.RedisManager // Redis instance for storing blacklisted tokens
}

// NewJWTFactory is a constructor function to create a new JWTFactory instance.
// It takes the JWTConfig and Redis instance as arguments and returns a new JWTFactory object.
func NewJWTFactory(c JWTConfig, r redis.RedisManager) *JWTFactory {
	return &JWTFactory{c, r}
}

// CreateJWT generates a new JWT token with the provided claims.
// It uses the HMAC signing method and the configured secret key to sign the token.
func (f *JWTFactory) CreateJWT(claims *JWTClaims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(f.cfg.Key))
}

// VerifyJWT validates the JWT token passed as a string and returns the claims if the token is valid.
// It checks if the token is blacklisted, validates the signing method, and verifies the token's integrity.
func (f *JWTFactory) VerifyJWT(ctx context.Context, tokenString string) (*JWTClaims, error) {
	if f.IsTokenBlacklisted(ctx, tokenString) {
		return nil, response.Unauthorized(app_error.ErrInvalidToken)
	}

	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key for validation
		return []byte(f.cfg.Key), nil
	})
	if err != nil {
		return nil, response.Unauthorized(app_error.ErrInvalidToken)
	}

	if !token.Valid {
		return nil, response.Unauthorized(app_error.ErrInvalidToken)
	}

	return token.Claims.(*JWTClaims), nil
}

// AddToBlacklist adds the JWT token to the blacklist with the specified expiration time.
// It stores the token in Redis to prevent its use in future requests.
func (f *JWTFactory) AddToBlacklist(context context.Context, token string, expiration time.Time) error {
	return f.redis.Set(context, token, "revoked", time.Until(expiration))
}

// IsTokenBlacklisted checks if the provided JWT token is blacklisted in Redis.
// If the token is found in Redis with the revoked value, it is considered blacklisted.
func (f *JWTFactory) IsTokenBlacklisted(ctx context.Context, token string) bool {
	revoked, err := f.redis.Get(ctx, token)
	if err != nil {
		return false
	}

	// Return true if the token is blacklisted, false otherwise
	return revoked == "revoked"
}
