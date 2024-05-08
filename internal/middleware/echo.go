package middleware

import (
	"net/http"

	"github.com/DoWithLogic/golang-clean-architecture/config"

	"github.com/DoWithLogic/golang-clean-architecture/pkg/observability"
	"github.com/go-playground/validator/v10"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (v *CustomValidator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

// configCORS contains the CORS (Cross-Origin Resource Sharing) configuration for the server.
var configCORS = echoMiddleware.CORSConfig{
	AllowOrigins: []string{"*"},
	AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch},
}

// NewEchoServer creates and configures a new Echo server instance.
// Parameters:
//   - cfg: The application configuration.
//
// Returns:
//   - *echo.Echo: A configured Echo server instance.
func NewEchoServer(cfg config.Config) *echo.Echo {
	e := echo.New()
	e.Use(echoMiddleware.RecoverWithConfig(echoMiddleware.RecoverConfig{DisableStackAll: true}))
	e.Use(echoMiddleware.CORSWithConfig(configCORS))
	e.Use(echoprometheus.NewMiddleware("http"))
	e.Use(LoggingMiddleware(observability.NewZeroLogHook().Z()))
	e.Use(CacheWithRevalidation)
	e.Validator = &CustomValidator{validator: validator.New()}

	if cfg.Observability.Enable {
		e.Use(otelecho.Middleware(cfg.App.Name))
	}

	e.HTTPErrorHandler = errorHandler
	e.Debug = cfg.Server.Debug

	return e
}
