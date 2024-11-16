package app

import (
	"context"
	"net/http"

	"github.com/DoWithLogic/golang-clean-architecture/internal/middleware"
	userV1 "github.com/DoWithLogic/golang-clean-architecture/internal/users/delivery/http/v1"
	userRepository "github.com/DoWithLogic/golang-clean-architecture/internal/users/repository"
	userUseCase "github.com/DoWithLogic/golang-clean-architecture/internal/users/usecase"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/app_crypto"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/app_jwt"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/app_redis"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

func (app *App) startService() error {
	domain := app.echo.Group("/api/v1")
	domain.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello Word 👋")
	})

	client := redis.NewClient(&redis.Options{
		Addr:     app.cfg.Redis.Addr,
		Password: app.cfg.Redis.Password,
		DB:       app.cfg.Redis.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return err
	}

	var (
		redis      = app_redis.NewRedis(client)
		crypto     = app_crypto.NewCrypto(app.cfg.Authentication.Key)
		jwt        = app_jwt.NewJWT(app.cfg.JWT, redis)
		middleware = middleware.NewMiddleware(jwt)

		userRepo = userRepository.NewRepository(app.db)
		userUC   = userUseCase.NewUseCase(userRepo, jwt, crypto)
		userCTRL = userV1.NewHandlers(userUC)
	)

	userCTRL.UserRoutes(domain, middleware)

	return nil
}
