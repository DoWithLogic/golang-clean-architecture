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

		// Create a context with a timeout to ensure graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Attempt to gracefully shutdown the Echo server
		if err := app.echo.Shutdown(ctx); err != nil {
			log.Errorf("Server shutdown failed: %v", err)
		} else {
			log.Info("Server shutdown gracefully")
		}

		// Close the database connection
		if err := app.db.DB.Close(); err != nil {
			log.Errorf("Database connection close failed: %v", err)
		} else {
			log.Info("Database connection closed")
		}
	}()

	// Start the Echo server
	port := fmt.Sprintf(":%s", app.cfg.Server.Port)
	log.Infof("Server is starting on port %s", port)

	if err := app.echo.Start(port); err != nil {
		log.Errorf("Error starting server: %v", err)
		return err
	}

	return nil
}
