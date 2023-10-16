package app

import (
	userV1 "github.com/DoWithLogic/golang-clean-architecture/internal/users/delivery/http/v1"
	userRepository "github.com/DoWithLogic/golang-clean-architecture/internal/users/repository"
	userUseCase "github.com/DoWithLogic/golang-clean-architecture/internal/users/usecase"
)

func (app *App) StartService() error {
	// define repository
	userRepo := userRepository.NewRepository(app.DB, app.Log)

	// define usecase
	userUC := userUseCase.NewUseCase(userRepo, app.Log, app.Cfg)

	// define controllers
	userCTRL := userV1.NewHandlers(userUC, app.Log)

	version := app.Echo.Group("/api/v1/")

	userV1.UserPrivateRoute(version, userCTRL, app.Cfg)

	return nil
}
