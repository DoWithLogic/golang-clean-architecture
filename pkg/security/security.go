package security

import (
	"context"
	"errors"
	"fmt"

	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/response"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/response/app_error"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/types"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type (
	PayloadToken struct {
		Data *Data `json:"data"`
		jwt.StandardClaims
	}

	Data struct {
		ID           int64              `json:"id"`
		ContactType  types.CONTACT_TYPE `json:"contact_type"`
		ContactValue string             `json:"contact_value"`
	}
)

var (
	ErrInvalidToken = errors.New("invalid authentication token")
)

type JWT struct {
	cfg config.JWTConfig
}

func NewJWT(cfg config.JWTConfig) *JWT {
	return &JWT{cfg: cfg}
}

func (j *JWT) GenerateToken(ctx context.Context, request PayloadToken) (token string, err error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, request).SignedString([]byte(j.cfg.Key))
}

func (j *JWT) ValidateToken(c echo.Context, token string) error {
	tokenWithoutBearer := token[len("Bearer "):]

	newToken, err := jwt.ParseWithClaims(tokenWithoutBearer, &PayloadToken{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.cfg.Key), nil
	})

	if err != nil {
		return response.ErrorBuilder(response.Unauthorized(ErrInvalidToken)).Send(c)
	}
	if !newToken.Valid {
		return response.ErrorBuilder(response.Unauthorized(ErrInvalidToken)).Send(c)
	}

	// Store the token claims in the request context for later use
	c.Set(types.CredentialDataContextKey.String(), newToken.Claims.(*PayloadToken))

	return nil
}

func NewTokenInformation(ctx echo.Context) (*PayloadToken, error) {
	tokenInformation, ok := ctx.Get(types.CredentialDataContextKey.String()).(*PayloadToken)
	if !ok {
		return tokenInformation, response.Unauthorized(app_error.ErrFailedGetTokenInformation)
	}

	return tokenInformation, nil
}
