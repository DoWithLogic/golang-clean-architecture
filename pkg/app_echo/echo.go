package app_echo

import (
	"net/http"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"

	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

type EchoConfig struct {
	Port  string // The port on which the server will listen.
	Debug bool   // Indicates if debug mode is enabled.
}

type CORSConfig echoMiddleware.CORSConfig

type echoRequest struct {
	IsObservabilityEnable bool
	ServiceName           *string
	CORSConfig            *CORSConfig
}

type EchoOptionFn func(*echoRequest)

func WithTracing(serviceName string) EchoOptionFn {
	return func(e *echoRequest) {
		e.IsObservabilityEnable = true
		e.ServiceName = &serviceName
	}
}
func WithCORS(c CORSConfig) EchoOptionFn { return func(er *echoRequest) { er.CORSConfig = &c } }

var defaultCORSConfig = CORSConfig{
	AllowOrigins: []string{"*"},
	AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
}

func defaultEchoRequest() *echoRequest {
	return &echoRequest{
		CORSConfig:            &defaultCORSConfig,
		IsObservabilityEnable: false,
	}
}

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
