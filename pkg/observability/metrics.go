package observability

import (
	"context"
	"time"

	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// OTLP communication modes
const (
	OTLP_HTTP_MODE = "otlp/http" // HTTP mode for OTLP communication
	OTLP_GRPC_MODE = "otlp/grpc" // gRPC mode for OTLP communication
	CONSOLE_MODE   = "console"   // Console mode, possibly for debugging
)

// InitMeterProvider initializes and returns an OpenTelemetry MeterProvider based on the provided configuration.
// The function sets up a metric exporter based on the observability mode specified in the configuration and creates a MeterProvider with the specified interval for exporting metrics.
//
// Parameters:
// - config: Configuration containing observability mode and application details.
//
// Returns:
// - A pointer to an OpenTelemetry MeterProvider.
// - An error if any occurs during initialization or setting up the provider.
func InitMeterProvider(observability config.ObservabilityConfig, app config.AppConfig) (*sdkmetric.MeterProvider, error) {
	var (
		exporter sdkmetric.Exporter
		err      error
	)

	// Determine the type of metric exporter based on the observability mode.
	switch observability.Mode {
	case OTLP_HTTP_MODE:
		exporter, err = otlpmetrichttp.New(
			context.Background(),
			otlpmetrichttp.WithInsecure(),
		)
	case OTLP_GRPC_MODE:
		exporter, err = otlpmetricgrpc.New(
			context.Background(),
			otlpmetricgrpc.WithInsecure(),
		)
	default:
		return nil, errors.New("invalid observability mode")
	}

	if err != nil {
		return nil, err
	}

	// Create and configure the MeterProvider with a periodic reader for exporting metrics.
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(2*time.Second))),
		sdkmetric.WithResource(initResource(app.Name, app.Version, app.Environment)),
	)

	// Set the MeterProvider as the global meter provider.
	otel.SetMeterProvider(mp)

	return mp, nil
}
