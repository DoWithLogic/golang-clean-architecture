package v1

import (
	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/middleware"
	"github.com/labstack/echo/v4"
)

func UserPrivateRoute(version *echo.Group, h Handlers, cfg config.Config) {
	users := version.Group("users")
	users.POST("", h.CreateUser)
	users.PATCH("/:id", h.UpdateUser, middleware.AuthorizeJWT(cfg))
	users.PUT("/:id", h.UpdateUserStatus, middleware.AuthorizeJWT(cfg))
}
