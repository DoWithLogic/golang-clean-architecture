package app

import (
	"context"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/datasource"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/middleware"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/otel/zerolog"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type App struct {
	db   *sqlx.DB
	echo *echo.Echo
	log  *zerolog.Logger
	cfg  config.Config
}

func NewApp(ctx context.Context, cfg config.Config) *App {
	db, err := datasource.NewDatabase(cfg.Database)
	if err != nil {
		panic(err)
	}

	return &App{
		db:   db,
		echo: echo.New(),
		log:  zerolog.NewZeroLog(ctx, os.Stdout),
		cfg:  cfg,
	}
}

func (app *App) Start() error {
	if err := app.StartService(); err != nil {
		app.log.Z().Err(err).Msg("[app]StartService")

		return err
	}

	app.echo.Debug = app.cfg.Server.Debug
	app.echo.Use(middleware.AppCors())
	app.echo.Use(middleware.CacheWithRevalidation)

	return app.echo.Start(fmt.Sprintf(":%s", app.cfg.Server.RESTPort))
}
