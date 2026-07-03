// Package logging provides a reusable Echo middleware for structured HTTP
// request/response logging via zerolog. It supports sensitive-field masking,
// path exclusion, and X-Request-ID propagation.
//
// # Quick Start
//
//	import (
//	    logging "github.com/golang-clean-architecture/pkg/logging"
//	)
//
//	e := echo.New()
//	e.Use(logging.Middleware(
//	    logging.WithLogger(&logger),
//	    logging.WithMaskedKeys("password", "token", "secret"),
//	))
package logging

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

// loggingContext bundles all data required to emit a single structured log
// entry for an inbound HTTP request and its response.
type loggingContext struct {
	Request         *http.Request
	RequestBodyJSON map[string]any
	StartTime       time.Time
	Options         loggingConfig
	Logger          *zerolog.Logger
	Ctx             context.Context
}
