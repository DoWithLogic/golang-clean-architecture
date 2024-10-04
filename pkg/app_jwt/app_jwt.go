package app_jwt

import (
	"context"

	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/DoWithLogic/golang-clean-architecture/internal/middleware"
	"github.com/golang-jwt/jwt"
)

type JWT struct {
	cfg config.JWTConfig
}

func NewJWT(cfg config.JWTConfig) *JWT {
	return &JWT{cfg: cfg}
}

func (j *JWT) GenerateToken(ctx context.Context, request middleware.PayloadToken) (token string, err error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, request).SignedString([]byte(j.cfg.Key))
}
