package v1

import (
	"github.com/DoWithLogic/golang-clean-architecture/internal/middleware"
	"github.com/labstack/echo/v4"
)

func (h *handlers) UserRoutes(domain *echo.Group, mw *middleware.Middleware) {
	// privates
	{
		private := domain.Group("/user", mw.JWTMiddleware())
		private.GET("/detail", h.UserDetail)
		private.PATCH("/update", h.UpdateUser)
		private.PUT("/update/status", h.UpdateUserStatus)
	}

	// publics
	domain.POST("/user", h.CreateUser)
	domain.POST("/user/login", h.Login)
}
