package observability

import (
	"context"
	"io"
	"log"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// Logger: A utility structure containing both a standard Go logger and a zerolog logger.
// It also provides methods to retrieve each logger instance and to set the log level dynamically.
type (
	Logger struct {
		standard *log.Logger     // standard: Standard Go logger.
		zerolog  *zerolog.Logger // zerolog: Zerolog logger.
	}
)

// TracingHook: Custom hook for zerolog that extracts tracing information from the context and adds it to the log entry.
type TracingHook struct{}

// Run: Implements the Run method of the zerolog.Hook interface to customize log events with tracing information.
func (h TracingHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	ctx := e.GetCtx()
	span := trace.SpanContextFromContext(ctx)

	if span.HasTraceID() {
		e.Str("span_id", span.SpanID().String()).Str("trace_id", span.TraceID().String())
	}

}

// NewZeroLogHook: Creates a new logger with a zerolog hook for tracing.
func NewZeroLogHook() *Logger {
	z := zerolog.New(os.Stdout).Hook(TracingHook{}).With().Timestamp().Stack().Logger()

	return &Logger{log.New(z, "", 0), &z}
}

// NewZeroLog: Creates a new logger with tracing information extracted from the provided context.
func NewZeroLog(ctx context.Context, c ...io.Writer) *Logger {
	span := trace.SpanContextFromContext(ctx)

	z := zerolog.New(os.Stdout).With().Timestamp().
		Str("span_id", span.SpanID().String()).
		Str("trace_id", span.TraceID().String()).
		Stack().Logger()

	return &Logger{log.New(z, "", 0), &z}
}

// S: Returns the standard Go logger from the Logger structure.
func (x *Logger) S() *log.Logger { return x.standard }

// Z: Returns the zerolog logger from the Logger structure.
func (x *Logger) Z() *zerolog.Logger { return x.zerolog }

// Level: Sets the log level for the zerolog logger and updates the standard logger accordingly.
func (x *Logger) Level(level string) *Logger {
	lv, err := zerolog.ParseLevel(strings.ToLower(level))
	if err == nil {
		*x.zerolog = x.zerolog.Level(lv)
		x.standard.SetOutput(x.zerolog)
	}

	return x
}
