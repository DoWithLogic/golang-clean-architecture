package app

import (
	"net/http"

	userV1 "github.com/DoWithLogic/golang-clean-architecture/internal/users/delivery/http/v1"
	userRepository "github.com/DoWithLogic/golang-clean-architecture/internal/users/repository"
	userUseCase "github.com/DoWithLogic/golang-clean-architecture/internal/users/usecase"
	"github.com/labstack/echo/v4"
)

func (app *App) startService() error {
	userRepo := userRepository.NewRepository(app.db)
	userUC := userUseCase.NewUseCase(userRepo, app.cfg)
	userCTRL := userV1.NewHandlers(userUC)

	domain := app.echo.Group("/api/v1/users")
	domain.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello Word 👋")
	})

	userCTRL.UserRoutes(domain, app.cfg)

	return nil
}
