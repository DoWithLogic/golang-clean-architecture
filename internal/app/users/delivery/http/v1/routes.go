package v1

import (
	"github.com/DoWithLogic/golang-clean-architecture/internal/middleware"
	"github.com/labstack/echo/v4"
)

func (h *handlers) MapRoutes(echo *echo.Group, mw *middleware.Middleware) {
	h.registerPublicRoutes(echo.Group("/user/public"))
	h.registerPrivateRoutes(echo.Group("/user", mw.JWTMiddleware()))
}

func (h *handlers) registerPublicRoutes(echo *echo.Group) {
	echo.POST("/login", h.LoginHandler)
	echo.POST("/sign-up", h.SignUpHandler)
}

func (h *handlers) registerPrivateRoutes(echo *echo.Group) {
	echo.GET("/:id/detail", h.UserDetailByIDHandler)
	echo.GET("/contact/:contact_value/detail", h.UserDetailByContactValueHandler)
	echo.PATCH("/:id/update", h.UpdateUserHandler)
	echo.PUT("/:id/status/transition", h.TransitionUserStatusHandler)
}
