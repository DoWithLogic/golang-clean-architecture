package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

func LoggingMiddleware(logger *zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			// Track the start time of the request
			startTime := time.Now()

			if errCheck := func() error {
				// Check if the request has a body
				if req.Body != nil {
					// Read the request body
					reqBody, err := io.ReadAll(req.Body)
					if err != nil {
						return err
					}

					// Restore the request body so it can be used in subsequent handlers
					req.Body = io.NopCloser(bytes.NewBuffer(reqBody))

					// Include request body as JSON in the log
					var requestBodyJSON interface{}
					if err := json.Unmarshal(reqBody, &requestBodyJSON); err != nil {
						return err
					}

					// Include request body in the log
					logger.Info().
						Str("method", req.Method).
						Str("path", req.URL.Path).
						Str("ip", req.RemoteAddr).
						Interface("body", requestBodyJSON).
						Msg("Received request")
				}

				return nil
			}(); errCheck != nil {
				// Log trace and span IDs using zerolog
				logger.Info().
					Str("method", req.Method).
					Str("path", req.URL.Path).
					Str("ip", req.RemoteAddr).
					Msg("Received request")
			}

			// Continue to the next middleware/handler
			err := next(c)

			// Log response details using zerolog
			logger.Info().
				Int("status", res.Status).
				Int64("size", res.Size).
				Dur("duration", time.Since(startTime)).
				Msg("Sent response")

			return err
		}
	}
}
