package app

import (
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/delivery/http/handler"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/delivery/http/route"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/repository"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/usecase"
)

func (app *App) StartService() error {
	// define repository
	repository := repository.NewRepository(app.DB, app.Log)

	// define usecase
	usecase := usecase.NewUseCase(repository, app.DB, app.Log)

	// define controllers
	controller := handler.NewHandlers(usecase, app.Log)

	version := app.Echo.Group("/api/v1/")

	route.RouteUsers(version, controller)

	return nil
}
