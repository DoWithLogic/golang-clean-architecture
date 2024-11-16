package app_jwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/app_redis"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/apperror"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/constant"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/response"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type (
	PayloadToken struct {
		Data *Data `json:"data"`
		jwt.StandardClaims
	}

	Data struct {
		UserID int64  `json:"user_id"`
		Email  string `json:"email"`
	}
)

var (
	errMissingJwtToken = errors.New("Missing JWT token")
	errInvalidJwtToken = errors.New("Invalid JWT token")
)

type JWT struct {
	cfg   config.JWTConfig
	redis app_redis.Redis
}

func NewJWT(cfg config.JWTConfig, redis app_redis.Redis) *JWT {
	return &JWT{cfg: cfg}
}

func (j *JWT) GenerateToken(ctx context.Context, request PayloadToken) (token string, err error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, request).SignedString([]byte(j.cfg.Key))
}

func (j *JWT) ValidateToken(c echo.Context, token string) error {
	tokenWithoutBearer := token[len("Bearer "):]
	if j.IsTokenRevoked(c.Request().Context(), tokenWithoutBearer) {
		return response.ErrorBuilder(apperror.Unauthorized(errInvalidJwtToken)).Send(c)
	}

	newToken, err := jwt.ParseWithClaims(tokenWithoutBearer, &PayloadToken{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.cfg.Key), nil
	})

	if err != nil {
		return response.ErrorBuilder(apperror.Unauthorized(errInvalidJwtToken)).Send(c)
	}
	if !newToken.Valid {
		return response.ErrorBuilder(apperror.Unauthorized(errInvalidJwtToken)).Send(c)
	}

	// Store the token claims in the request context for later use
	c.Set(constant.AuthCredentialKey, newToken.Claims.(*PayloadToken))

	return nil
}

func (j *JWT) RevokeToken(ctx context.Context, token string, expiration time.Duration) error {
	return j.redis.Set(ctx, token, constant.TokenRevoked, expiration)
}

func (j *JWT) IsTokenRevoked(ctx context.Context, token string) bool {
	revoked, err := j.redis.Get(ctx, token)
	if err != nil {
		return false
	}

	return revoked == constant.TokenRevoked
}

func NewTokenInformation(ctx echo.Context) (*PayloadToken, error) {
	tokenInformation, ok := ctx.Get(constant.AuthCredentialKey).(*PayloadToken)
	if !ok {
		return tokenInformation, apperror.Unauthorized(apperror.ErrFailedGetTokenInformation)
	}

	return tokenInformation, nil
}
