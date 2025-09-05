package middleware

import (
	"github.com/DoWithLogic/golang-clean-architecture/pkg/app_jwt"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/apperror"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/constants"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/response"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

var (
	errMissingJwtToken = errors.New("Missing JWT token")
	errInvalidJwtToken = errors.New("Invalid JWT token")
)

type Middleware struct {
	jwt *app_jwt.JWT
}

func NewMiddleware(jwt *app_jwt.JWT) *Middleware {
	return &Middleware{jwt: jwt}
}

// Middleware function to validate JWT token
func (m *Middleware) JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get(constants.AuthorizationHeaderKey)
			if token == "" {
				return response.ErrorBuilder(apperror.Unauthorized(errMissingJwtToken)).Send(c)
			}

			if err := m.jwt.ValidateToken(c, token); err != nil {
				return err
			}

			return next(c)
		}
	}
}
