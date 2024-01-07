package v1

import (
	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/middleware"
	"github.com/labstack/echo/v4"
)

func (h *handlers) UserRoutes(domain *echo.Group, cfg config.Config) {
	domain.POST("/", h.CreateUser)
	domain.POST("/login", h.Login)
	domain.GET("/detail", h.UserDetail, middleware.AuthorizeJWT(cfg))
	domain.PATCH("/update", h.UpdateUser, middleware.AuthorizeJWT(cfg))
	domain.PUT("/update/status", h.UpdateUserStatus, middleware.AuthorizeJWT(cfg))
}
