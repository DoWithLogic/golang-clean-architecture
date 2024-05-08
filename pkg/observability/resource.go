package observability

import (
	"context"

	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// initResource initializes and returns an OpenTelemetry resource based on the provided application details.
// This function creates a composite resource that combines default resource attributes with additional attributes such as OS, process, container, host, service name, service version,
// and deployment environment.
//
// Parameters:
// - appName: Name of the application.
// - appVersion: Version of the application.
// - appEnv: Environment (e.g., development, production) of the application.
//
// Returns:
// - A pointer to an OpenTelemetry resource with the combined attributes.
func initResource(appName string, appVersion string, appEnv string) *sdkresource.Resource {
	// Create additional resource attributes based on the application details.
	extraResource, _ := sdkresource.New(
		context.Background(),
		sdkresource.WithOS(),
		sdkresource.WithProcess(),
		sdkresource.WithContainer(),
		sdkresource.WithHost(),
		sdkresource.WithAttributes(
			semconv.ServiceName(appName),
			semconv.ServiceVersion(appVersion),
			semconv.DeploymentEnvironment(appEnv),
		),
	)

	// Merge the additional resource attributes with the default attributes.
	resource, _ := sdkresource.Merge(
		sdkresource.Default(),
		extraResource,
	)

	return resource
}
