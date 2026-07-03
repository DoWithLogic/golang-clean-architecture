package middleware

import (
	"github.com/DoWithLogic/golang-clean-architecture/pkg/errs"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/response"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/security"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/types"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

var (
	errMissingJwtToken = errors.New("Missing JWT token")
)

type Middleware struct {
	jwt *security.JWT
}

func NewMiddleware(jwt *security.JWT) *Middleware {
	return &Middleware{jwt: jwt}
}

// Middleware function to validate JWT token
func (m *Middleware) JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get(types.AuthorizationHeaderKey.String())
			if token == "" {
				return response.ErrorBuilder(errs.Unauthorized(errMissingJwtToken)).Send(c)
			}

			if err := m.jwt.ValidateToken(c, token); err != nil {
				return err
			}

			return next(c)
		}
	}
}
