package server

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	userV1 "github.com/DoWithLogic/golang-clean-architecture/internal/app/users/delivery/http/v1"
	userRepository "github.com/DoWithLogic/golang-clean-architecture/internal/app/users/repository"
	userUseCase "github.com/DoWithLogic/golang-clean-architecture/internal/app/users/usecase"
	"github.com/DoWithLogic/golang-clean-architecture/internal/middleware"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/app_crypto"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/app_jwt"
	"github.com/labstack/echo/v4"

	echoSwagger "github.com/swaggo/echo-swagger"
)

func (s *Server) setup() error {
	domain := s.echo.Group("/api/v1")
	domain.GET("/ping", func(c echo.Context) error { return c.String(http.StatusOK, "Hello Word 👋") })

	// Serve Swagger documentation
	domain.GET("/swagger/*", echoSwagger.WrapHandler)

	customSwaggerHandler := func(c echo.Context) error {
		// Read the generated Swagger JSON file
		data, err := os.ReadFile("docs/swagger.json") // Adjust the path as necessary
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error reading Swagger JSON")
		}

		// Unmarshal the Swagger JSON into a map
		var swaggerDoc map[string]interface{}
		if err := json.Unmarshal(data, &swaggerDoc); err != nil {
			return c.String(http.StatusInternalServerError, "Error parsing Swagger JSON")
		}

		// Modify SwaggerInfo
		swaggerDoc["host"] = s.cfg.App.Host
		swaggerDoc["schemes"] = strings.Split(s.cfg.App.Scheme, ",")

		// Marshal the modified Swagger JSON back to a string
		modifiedSwaggerJSON, err := json.Marshal(swaggerDoc)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error generating modified Swagger JSON")
		}

		return c.String(http.StatusOK, string(modifiedSwaggerJSON))
	}
	domain.GET("/swagger/doc.json", customSwaggerHandler)

	var (
		crypto     = app_crypto.NewCrypto(s.cfg.Authentication.Key)
		jwt        = app_jwt.NewJWT(s.cfg.JWT)
		middleware = middleware.NewMiddleware(jwt)

		userRepo = userRepository.NewRepository(s.db)
		userUC   = userUseCase.NewUseCase(userUseCase.Dependencies{
			UseCases:     userUseCase.UseCases{},
			Repositories: userUseCase.Repositories{Repo: userRepo},
			Pkgs:         userUseCase.Pkgs{AppJwt: jwt, Crypto: crypto},
		})
		userHandler = userV1.NewHandlers(userUC)
	)

	handlers := []routeMapper{
		userHandler,
	}

	// Routes
	for idx := range handlers {
		handlers[idx].MapRoutes(domain, middleware)
	}

	return nil
}

type routeMapper interface {
	MapRoutes(echo *echo.Group, mw *middleware.Middleware)
}
