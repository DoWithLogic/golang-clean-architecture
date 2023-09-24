package route

import (
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/delivery/http/handler"
	"github.com/labstack/echo/v4"
)

func RouteUsers(version *echo.Group, handler handler.Handlers) {
	users := version.Group("users")
	users.POST("", handler.CreateUser)
	users.PATCH("/:id", handler.UpdateUser)
	users.PUT("/:id", handler.UpdateUserStatus)
}
