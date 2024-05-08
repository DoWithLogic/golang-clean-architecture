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
	"github.com/DoWithLogic/golang-clean-architecture/internal/middleware"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/datasource"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"github.com/samber/lo"
)

type App struct {
	db   *sqlx.DB      // Database connection.
	echo *echo.Echo    // Echo HTTP server instance.
	cfg  config.Config // Configuration settings for the application.
}

func NewApp(ctx context.Context, cfg config.Config) *App {
	return &App{
		db:   lo.Must(datasource.NewDatabase(cfg.Database)),
		echo: middleware.NewEchoServer(cfg),
		cfg:  cfg,
	}
}

func (app *App) Run() error {
	if err := app.startService(); err != nil {
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

		// Create a context with a timeout of 10 seconds for the server shutdown.
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Shutdown gracefully.
		app.db.DB.Close()
		app.echo.Shutdown(ctx)
	}()

	return app.echo.Start(fmt.Sprintf(":%s", app.cfg.Server.Port))
}
