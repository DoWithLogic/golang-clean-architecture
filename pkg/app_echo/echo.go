package app_echo

import (
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"

	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

// NewEchoServer creates and configures a new Echo server instance.
// Parameters:
//   - cfg: The application configuration.
//
// Returns:
//   - *echo.Echo: A configured Echo server instance.
func (cfg EchoConfig) New(opts ...EchoOptionFn) *echo.Echo {
	request := defaultEchoRequest()
	for _, opt := range opts {
		opt(request)
	}

	e := echo.New()
	e.Use(echoMiddleware.RecoverWithConfig(echoMiddleware.RecoverConfig{DisableStackAll: true}))
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig(*request.CORSConfig)))
	e.Use(echoprometheus.NewMiddleware("http"))
	e.Use(cacheWithRevalidation)

	if cfg.Debug {
		e.Debug = cfg.Debug
	}

	if request.IsObservabilityEnable {
		e.Use(otelecho.Middleware(*request.ServiceName))
	}

	e.HTTPErrorHandler = errorHandler

	return e
}
