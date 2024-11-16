package observability

import (
	"context"

	"github.com/DoWithLogic/golang-clean-architecture/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// initializeTracerProvider initializes and configures a tracer provider based on the observability mode from the configuration.
// It supports OTLP gRPC, OTLP HTTP, and stdout (default) exporters for trace data.
// Returns a tracer provider instance or an error if initialization fails.
func InitTracerProvider(observabilityCfg config.ObservabilityConfig, appCfg config.AppConfig) (*sdktrace.TracerProvider, error) {
	var (
		exporter sdktrace.SpanExporter
		err      error
	)

	// Determine the exporter type based on the observability mode.
	switch observabilityCfg.Mode {
	case OTLP_GRPC_MODE:
		exporter, err = otlptracegrpc.New(
			context.Background(),
			otlptracegrpc.WithInsecure(),
		)
	case OTLP_HTTP_MODE:
		exporter, err = otlptracehttp.New(
			context.Background(),
			otlptracehttp.WithInsecure(),
		)
	default:
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
	}

	if err != nil {
		return nil, err
	}

	// Configure and set up the tracer provider with the chosen exporter and other necessary configurations.
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(initResource(appCfg.Name, appCfg.Version, appCfg.Environment)),
	)

	// Set the global tracer provider for the OpenTelemetry API.
	otel.SetTracerProvider(tracerProvider)
	// Set the text map propagator for trace context and baggage.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tracerProvider, nil
}
