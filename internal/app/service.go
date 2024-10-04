package app

import (
	"net/http"

	userV1 "github.com/DoWithLogic/golang-clean-architecture/internal/users/delivery/http/v1"
	userRepository "github.com/DoWithLogic/golang-clean-architecture/internal/users/repository"
	userUseCase "github.com/DoWithLogic/golang-clean-architecture/internal/users/usecase"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/app_crypto"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/app_jwt"
	"github.com/labstack/echo/v4"
)

func (app *App) startService() error {
	domain := app.echo.Group("/api/v1/users")
	domain.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello Word 👋")
	})

	var (
		crypto = app_crypto.NewCrypto(app.cfg.Authentication.Key)
		appJwt = app_jwt.NewJWT(app.cfg.JWT)

		userRepo = userRepository.NewRepository(app.db)
		userUC   = userUseCase.NewUseCase(userRepo, appJwt, crypto)
		userCTRL = userV1.NewHandlers(userUC)
	)

	return userCTRL.UserRoutes(domain, app.cfg)
}
