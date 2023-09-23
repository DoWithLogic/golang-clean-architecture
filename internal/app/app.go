package app

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/DoWithLogic/golang-clean-architecture/internal/infrastructure"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/otel/zerolog"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type App struct {
	DB   *sqlx.DB
	Echo *echo.Echo
	Log  *zerolog.Logger
	Cfg  config.Config
}

func NewApp(ctx context.Context, cfg config.Config) *App {
	db, err := infrastructure.NewDatabase(cfg.Database)
	if err != nil {
		panic(err)
	}

	echo, err := infrastructure.NewEcho(cfg.Server)
	if err != nil {
		panic(err)
	}

	return &App{
		DB:   db,
		Echo: echo,
		Log:  zerolog.NewZeroLog(ctx, os.Stdout),
		Cfg:  cfg,
	}
}

func (app *App) Start() error {
	if err := app.StartService(); err != nil {
		app.Log.Z().Err(err).Msg("[app]StartService")

		return err
	}

	return app.Echo.StartServer(&http.Server{
		Addr:         fmt.Sprintf(":%s", app.Cfg.Server.RESTPort),
		ReadTimeout:  app.Cfg.Server.ReadTimeout,
		WriteTimeout: app.Cfg.Server.WriteTimeout,
	})
}
