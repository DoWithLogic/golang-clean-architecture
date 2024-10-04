package v1

import (
	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/DoWithLogic/golang-clean-architecture/internal/middleware"
	"github.com/labstack/echo/v4"
)

func (h *handlers) UserRoutes(domain *echo.Group, cfg config.Config) error {
	domain.POST("", h.CreateUser)
	domain.POST("/login", h.Login)
	domain.GET("/detail", h.UserDetail, middleware.JWTMiddleware(cfg))
	domain.PATCH("/update", h.UpdateUser, middleware.JWTMiddleware(cfg))
	domain.PUT("/update/status", h.UpdateUserStatus, middleware.JWTMiddleware(cfg))

	return nil
}
