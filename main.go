package main

import (
	"context"
	"time"

	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/DoWithLogic/golang-clean-architecture/internal/server"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/observability"
	"github.com/labstack/gommon/log"
)

// main is entrypoint of application
//
//	@title						golang-clean-architecture
//	@version					v0.0.1
//	@description.markdown
//	@termsOfService				https://github.com/DoWithLogic/golang-clean-architecture
//	@BasePath					/api/v1
//
//	@securityDefinitions.apiKey	BearerToken
//	@in							header
//	@name						Authorization
func main() {
	// Load the application configuration from the specified directory.
	cfg, err := config.LoadConfig("config")
	if err != nil {
		// If an error occurs while loading the configuration, panic with the error.
		panic(err)
	}

	// Set the time zone to the specified value from the configuration.
	_, err = time.LoadLocation(cfg.Server.TimeZone)
	if err != nil {
		// If an error occurs while setting the time zone, log the error and exit the function.
		log.Error("Error on setting the time zone: ", err)
		return
	}

	// Initialize observability components if observability is enabled in the configuration.
	if cfg.Observability.Enable {
		// Initialize the tracer provider for distributed tracing.
		tracer, err := observability.InitTracerProvider(cfg.Observability, cfg.App)
		if err != nil {
			log.Warn("Failed to initialize tracer: ", err)
		}

		// Initialize the meter provider for metrics collection.
		meter, err := observability.InitMeterProvider(cfg.Observability, cfg.App)
		if err != nil {
			log.Warn("Failed to initialize meter: ", err)
		}

		// Ensure that the tracer and meter are shut down when the main function exits.
		defer func() {
			if tracer != nil {
				tracer.Shutdown(context.Background())
			}
			if meter != nil {
				meter.Shutdown(context.Background())
			}
		}()
	}

	if err := server.NewServer(context.Background(), cfg).Run(); err != nil {
		panic(err)
	}
}
