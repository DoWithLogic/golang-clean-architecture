package middleware

import (
	"fmt"

	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/apperror"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/constant"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/response"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type PayloadToken struct {
	Data *Data `json:"data"`
	jwt.StandardClaims
}

type Data struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
}

// Middleware function to validate JWT token
func JWTMiddleware(cfg config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := c.Request().Header.Get(constant.AuthorizationHeaderKey)
			if tokenString == "" {
				return response.ErrorBuilder(apperror.Unauthorized(errors.New("Missing JWT token"))).Send(c)
			}

			// Remove "Bearer " prefix from token string
			tokenString = tokenString[len("Bearer "):]

			token, err := jwt.ParseWithClaims(tokenString, &PayloadToken{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(cfg.JWT.Key), nil
			})

			if err != nil {
				return response.ErrorBuilder(apperror.Unauthorized(errors.New("Invalid JWT token"))).Send(c)
			}
			if !token.Valid {
				return response.ErrorBuilder(apperror.Unauthorized(errors.New("JWT token is not valid"))).Send(c)
			}

			// Store the token claims in the request context for later use
			claims := token.Claims.(*PayloadToken)
			c.Set(constant.AuthCredentialKey, claims)

			return next(c)
		}
	}
}

func NewTokenInformation(ctx echo.Context) (*PayloadToken, error) {
	tokenInformation, ok := ctx.Get(constant.AuthCredentialKey).(*PayloadToken)
	if !ok {
		return tokenInformation, apperror.Unauthorized(apperror.ErrFailedGetTokenInformation)
	}

	return tokenInformation, nil
}
