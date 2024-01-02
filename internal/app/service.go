package app

import (
	userV1 "github.com/DoWithLogic/golang-clean-architecture/internal/users/delivery/http/v1"
	userRepository "github.com/DoWithLogic/golang-clean-architecture/internal/users/repository"
	userUseCase "github.com/DoWithLogic/golang-clean-architecture/internal/users/usecase"
)

func (app *App) StartService() error {
	userRepo := userRepository.NewRepository(app.db)
	userUC := userUseCase.NewUseCase(userRepo, app.cfg)
	userCTRL := userV1.NewHandlers(userUC)

	domain := app.echo.Group("/api/v1/users")

	userCTRL.UserRoutes(domain, app.cfg)

	return nil
}
