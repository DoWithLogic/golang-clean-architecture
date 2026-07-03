package server

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	userV1 "github.com/DoWithLogic/golang-clean-architecture/internal/app/users/delivery/http/v1"
	userRepository "github.com/DoWithLogic/golang-clean-architecture/internal/app/users/repository"
	userUseCase "github.com/DoWithLogic/golang-clean-architecture/internal/app/users/usecase"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/encryptions"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/jwt"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/logging"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/middleware"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/observability"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/redis"
	"github.com/labstack/echo/v4"

	echoSwagger "github.com/swaggo/echo-swagger"
)

type routeMapper interface {
	MapRoutes(api *echo.Group, mw *middleware.Middleware)
}

func (s *Server) setup() error {
	s.setupMiddleware()

	api := s.echo.Group("/api/v1")

	s.registerUtilityRoutes(api)

	middleware, handlers := s.buildHandlers()

	for _, handler := range handlers {
		handler.MapRoutes(api, middleware)
	}

	return nil
}

func (s *Server) setupMiddleware() {
	logger := observability.NewZeroLogHook().Z()
	s.echo.Use(logging.Middleware(logging.WithLogger(logger), logging.WithMaskedKeys("password", "token")))
}

func (s *Server) registerUtilityRoutes(api *echo.Group) {
	api.GET("/ping", s.ping)
	api.GET("/swagger/*", echoSwagger.WrapHandler)
	api.GET("/swagger/doc.json", s.swaggerDoc)
}

func (s *Server) ping(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World 👋")
}

func (s *Server) swaggerDoc(c echo.Context) error {
	data, err := os.ReadFile("docs/swagger.json")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to read swagger document")
	}

	var doc map[string]any
	if err := json.Unmarshal(data, &doc); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to parse swagger document")
	}

	doc["host"] = s.cfg.App.Host
	doc["schemes"] = strings.Split(s.cfg.App.Scheme, ",")

	body, err := json.Marshal(doc)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate swagger document")
	}

	return c.Blob(http.StatusOK, echo.MIMEApplicationJSON, body)
}

func (s *Server) buildHandlers() (*middleware.Middleware, []routeMapper) {
	redisManager := redis.NewRedisManager(nil)

	jwtFactory := jwt.NewJWTFactory(s.cfg.JWT, redisManager)
	crypto := encryptions.NewCrypto(s.cfg.Authentication.Key)

	mw := middleware.New(jwtFactory)

	userRepo := userRepository.NewRepository(s.db)

	userUC := userUseCase.NewUseCase(userUseCase.Dependencies{
		Repositories: userUseCase.Repositories{
			Repo: userRepo,
		},
		Pkgs: userUseCase.Pkgs{
			AppJwt: jwtFactory,
			Crypto: crypto,
		},
	})

	handlers := []routeMapper{
		userV1.NewHandlers(userUC),
	}

	return mw, handlers
}
