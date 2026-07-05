package app_echo

import (
	"net/http"

	echoMiddleware "github.com/labstack/echo/v4/middleware"
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
