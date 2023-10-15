package v1

import (
	"github.com/labstack/echo/v4"
)

func UserPrivateRoute(version *echo.Group, h Handlers) {
	users := version.Group("users")
	users.POST("/", h.CreateUser)
	users.PATCH("/:id", h.UpdateUser)
	users.PUT("/:id", h.UpdateUserStatus)
}
