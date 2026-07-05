package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type SuccessDefault struct {
	Code    int             `json:"code" example:"200"` // HTTP status code.
	Message ResponseMessage `json:"message" example:"success"`
}

type SuccessResponse struct {
	SuccessDefault
	Data any `json:"data,omitempty"` // data payload.
}

// SuccessResponse represents a success response structure for API responses.
type SuccessWithPaginationResponse struct {
	SuccessResponse
	Meta any `json:"meta,omitempty"` //pagination payload.
}

// SuccessBuilder constructs a CustomResponse with a Success status and the provided response data.
func SuccessBuilder(response any, meta ...any) SuccessWithPaginationResponse {
	result := SuccessWithPaginationResponse{
		SuccessResponse: SuccessResponse{
			SuccessDefault: SuccessDefault{
				Code:    http.StatusOK,
				Message: SuccessMessage,
			},
			Data: response,
		},
	}

	if len(meta) > 0 {
		result.Meta = meta[0]
	}

	return result
}

// Send sends the CustomResponse as a JSON response using the provided Echo context.
func (c SuccessResponse) Send(ctx echo.Context) error {
	trace.SpanFromContext(ctx.Request().Context()).SetStatus(codes.Ok, http.StatusText(c.Code))
	return ctx.JSON(c.Code, c)
}
