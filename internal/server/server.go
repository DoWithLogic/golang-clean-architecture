package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/app_echo"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/datasources"
	appRedis "github.com/DoWithLogic/golang-clean-architecture/pkg/redis"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"github.com/samber/lo"

	"github.com/go-redis/redis/v8"
)

type Server struct {
	db          *gorm.DB      // Database connection.
	echo        *echo.Echo    // Echo HTTP server instance.
	cfg         config.Config // Configuration settings for the application.
	redisClient *redis.Client
}

func NewServer(ctx context.Context, cfg config.Config) *Server {
	serverOpts := []app_echo.EchoOptionFn{}
	if cfg.Observability.Enable {
		serverOpts = append(serverOpts, app_echo.WithTracing(cfg.App.Name))
	}

	return &Server{
		db:          lo.Must(datasources.NewMySQLDB(ctx, cfg.Database)),
		echo:        cfg.Server.New(serverOpts...),
		cfg:         cfg,
		redisClient: appRedis.NewRedisClient(ctx, cfg.Redis),
	}
}

func (s *Server) Run() error {
	if err := s.setup(); err != nil {
		return err
	}

	// Set up signal handling to gracefully shutdown the server upon receiving a SIGTERM or SIGINT signal.
	// Using a buffered channel with capacity 1 to ensure signals are not missed.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM)
	signal.Notify(quit, syscall.SIGINT)

	// Start a goroutine to handle the shutdown signal.
	go func() {
		// Wait for the signal from the quit channel.
		<-quit

		// Log the shutdown process.
		log.Info("Server is shutting down...")

		// Create a context with a timeout of 10 seconds for the server shutdown.
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Close DB Connection
		lo.Must(s.db.DB()).Close()

		// Shutdown gracefully.
		s.echo.Shutdown(ctx)
	}()

	// Start the echo server and listen on the configured port.
	return s.echo.Start(fmt.Sprintf(":%s", s.cfg.Server.Port))
}
