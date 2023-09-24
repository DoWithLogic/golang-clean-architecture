package route

import (
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/delivery/http/handler"
	"github.com/labstack/echo/v4"
)

func RouteUsers(version *echo.Group, ctrl handler.Handlers) {
	users := version.Group("users")
	users.POST("", ctrl.CreateUser)
	users.PATCH("/:id", ctrl.UpdateUser)
}
