package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/datasource"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/middleware"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/otel/zerolog"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
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

func (app *App) Run() error {
	if err := app.startService(); err != nil {
		app.log.Z().Err(err).Msg("[app]StartService")

		return err
	}

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM)
	signal.Notify(quit, syscall.SIGINT)

	go func() {
		<-quit
		log.Info("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		app.db.Close()
		app.echo.Shutdown(ctx)
	}()

	app.echo.Debug = app.cfg.Server.Debug
	app.echo.Use(middleware.AppCors())
	app.echo.Use(middleware.CacheWithRevalidation)

	return app.echo.Start(fmt.Sprintf(":%s", app.cfg.Server.RESTPort))
}
