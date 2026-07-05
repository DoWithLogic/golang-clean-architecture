package response

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ErrorResponse represents a failed response structure for API responses.
type ErrorResponse struct {
	Code    int             `json:"code" example:"500"`                      // HTTP status code.
	Message ResponseMessage `json:"message" example:"internal_server_error"` // Message corresponding to the status code.
	Error   string          `json:"error" example:"{$err}"`                  // error message.
}

// ErrorBuilder constructs a ErrorResponse based on the provided error.
func ErrorBuilder(err error) ErrorResponse {
	var appErr *AppError
	if errors.As(err, &appErr) {
		ae := err.(*AppError)

		return ErrorResponse{
			Code:    ae.Code,
			Message: ae.Message,
			Error:   ae.Error(),
		}
	}

	response := ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: InternalServerErrorMessage,
		Error:   InternalServerErrorMessage.String(),
	}

	if err != nil {
		response.Error = err.Error()
	}

	return response
}

// Send sends the CustomResponse as a JSON response using the provided Echo context.
func (x ErrorResponse) Send(c echo.Context) error {
	err := errors.New(x.Error)
	span := trace.SpanFromContext(c.Request().Context())
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())

	return c.JSON(x.Code, x)
}
