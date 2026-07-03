package response

import (
	"context"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// FailedResponse represents a failed response structure for API responses.
type FailedResponse struct {
	Code    int             `json:"code" example:"500"`                      // HTTP status code.
	Message ResponseMessage `json:"message" example:"internal_server_error"` // Message corresponding to the status code.
	Error   string          `json:"error" example:"{$err}"`                  // error message.
}

// BasicResponse represents a failed response structure for API responses.
type BasicResponse struct {
	Code    int    `json:"code" example:"500"`                      // HTTP status code.
	Message string `json:"message" example:"internal_server_error"` // Message corresponding to the status code.
	Error   string `json:"error" example:"{$err}"`                  // error message.
	Data    any    `json:"data,omitempty"`
}

// BasicBuilder constructs a BasicBuilder based on the provided error.
func BasicBuilder(result BasicResponse) BasicResponse {
	return result
}

// Send sends the CustomResponse as a JSON response using the provided Echo context.
func (c BasicResponse) Send(ctx echo.Context) error {
	trace.SpanFromContext(ctx.Request().Context()).SetStatus(codes.Ok, http.StatusText(c.Code))
	return ctx.JSON(c.Code, c)
}

// ErrorBuilder constructs a FailedResponse based on the provided error.
func ErrorBuilder(err error) FailedResponse {
	var appErr *AppError
	if errors.As(err, &appErr) {
		ae := err.(*AppError)

		return FailedResponse{
			Code:    ae.Code,
			Message: ae.Message,
			Error:   ae.Error(),
		}
	}

	var errStr = InternalServerErrorMessage.String()
	if err != nil {
		errStr = err.Error()
	}

	return FailedResponse{
		Code:    http.StatusInternalServerError,
		Message: InternalServerErrorMessage,
		Error:   errStr,
	}
}

// Send sends the CustomResponse as a JSON response using the provided Echo context.
func (x FailedResponse) Send(c echo.Context) error {
	err := errors.New(x.Error)
	span := trace.SpanFromContext(c.Request().Context())
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())

	return c.JSON(x.Code, x)
}

// SuccessResponse represents a success response structure for API responses.
type SuccessResponse struct {
	SuccessFormat
	Meta
}

type ResponseFormat struct {
	Code    int    `json:"code" example:"200"` // HTTP status code.
	Message string `json:"message" example:"success"`
}

type SuccessFormat struct {
	ResponseFormat
	Data any `json:"data,omitempty"` // data payload.
}

type Meta struct {
	Meta any `json:"meta,omitempty"` //pagination payload.
	SuccessFormat
}

// SuccessBuilder constructs a CustomResponse with a Success status and the provided response data.
func SuccessBuilder(response any, meta ...any) SuccessResponse {
	result := SuccessResponse{
		SuccessFormat: SuccessFormat{
			ResponseFormat: ResponseFormat{
				Code:    http.StatusOK,
				Message: SuccessMessage.String(),
			},
			Data: response,
		},
	}

	if len(meta) > 0 {
		result.Meta.Meta = meta[0]
	}

	return result
}

func SuccessNoResponseBuilder() SuccessResponse {
	result := SuccessResponse{
		SuccessFormat: SuccessFormat{
			ResponseFormat: ResponseFormat{
				Code:    http.StatusNoContent,
				Message: SuccessMessage.String(),
			},
		},
	}

	return result
}

// Send sends the CustomResponse as a JSON response using the provided Echo context.
func (c SuccessResponse) Send(ctx echo.Context) error {
	trace.SpanFromContext(ctx.Request().Context()).SetStatus(codes.Ok, http.StatusText(c.Code))
	return ctx.JSON(c.Code, c)
}

type validatable interface {
	Validate() error
}

type (
	handlerType any
)

func GenericHandler[Request validatable](c echo.Context, request *Request, usecaseFn handlerType) error {
	if err := c.Bind(request); err != nil {
		return ErrorBuilder(BadRequest(err)).Send(c)
	}

	if err := (*request).Validate(); err != nil {
		return ErrorBuilder(BadRequest(err)).Send(c)
	}

	fnValue := reflect.ValueOf(usecaseFn)
	fnType := fnValue.Type()

	if fnType.Kind() != reflect.Func {
		return ErrorBuilder(InternalServerError(errors.New("usecaseFn is not a function"))).Send(c)
	}

	// check input param
	if fnType.NumIn() != 2 {
		return ErrorBuilder(InternalServerError(errors.New("function must have 2 input parameters"))).Send(c)
	}

	// Check output number
	if fnType.NumOut() < 1 || fnType.NumOut() > 2 {
		return ErrorBuilder(InternalServerError(errors.New("function must have 1 or 2 return values"))).Send(c)
	}

	// Prepare arguments for the function call
	args := []reflect.Value{
		reflect.ValueOf(c.Request().Context()),
		reflect.ValueOf(*request),
	}

	// Call the function
	results := fnValue.Call(args)

	// Handle return values based on the number of results
	switch len(results) {
	case 1: // HandlerFn: returns only an error
		err := results[0].Interface()
		if err != nil {
			return ErrorBuilder(err.(error)).Send(c)
		}

		return SuccessBuilder(nil).Send(c)

	case 2: // HandlerWithResponseFn: returns a response and an error
		response := results[0].Interface()
		err := results[1].Interface()

		if err != nil {
			return ErrorBuilder(err.(error)).Send(c)
		}

		return SuccessBuilder(response).Send(c)

	default:
		return ErrorBuilder(InternalServerError(errors.New("unsupported function signature"))).Send(c)
	}
}

func GenericPaginationHandler[Request validatable, Response any](c echo.Context, request *Request, usecaseFn func(ctx context.Context, request *Request) (Response, error)) error {
	if err := c.Bind(request); err != nil {
		return ErrorBuilder(BadRequest(err)).Send(c)
	}

	if err := (*request).Validate(); err != nil {
		return ErrorBuilder(BadRequest(err)).Send(c)
	}

	result, err := usecaseFn(c.Request().Context(), request)
	if err != nil {
		return ErrorBuilder(err).Send(c)
	}

	return SuccessBuilder(result, *request).Send(c)
}
