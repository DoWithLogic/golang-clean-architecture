package middleware

import (
	"errors"
	"strings"

	"github.com/DoWithLogic/golang-clean-architecture/pkg/jwt"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/response"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/types"
	"github.com/labstack/echo/v4"
)

var (
	ErrInvalidAuthenticationCredentials = errors.New("invalid authentication credentials")
)

type Middleware struct {
	jwtFactory *jwt.JWTFactory
}

func New(jwtFactory *jwt.JWTFactory) *Middleware {
	return &Middleware{jwtFactory: jwtFactory}
}

func (m *Middleware) JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := c.Request().Header.Get(types.AuthorizationHeaderKey.String())
			if tokenString == "" {
				return response.ErrorBuilder(response.Unauthorized(ErrInvalidAuthenticationCredentials)).Send(c)
			}

			tokenString = strings.TrimPrefix(tokenString, "Bearer ")

			claims, err := m.jwtFactory.VerifyJWT(c.Request().Context(), tokenString)
			if err != nil {
				return response.ErrorBuilder(response.Unauthorized(ErrInvalidAuthenticationCredentials)).Send(c)
			}

			embedClaimedDataIntoContext(c, embedClaimedDataIntoContextOpts{claimedData: claims})

			return next(c)
		}
	}
}
