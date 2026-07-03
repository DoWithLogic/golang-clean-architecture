package response

import "net/http"

type AppError struct {
	Code    int
	Err     error
	Message ResponseMessage
}

func (e *AppError) Unwrap() error { return e.Err }
func (h AppError) Error() string  { return h.Err.Error() }

func BadRequest(err error) error {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: BadRequestMessage,
		Err:     err,
	}
}

func InternalServerError(err error) error {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: InternalServerErrorMessage,
		Err:     err,
	}
}

func Unauthorized(err error) error {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: UnauthorizedMessage,
		Err:     err,
	}
}

func Forbidden(err error) error {
	return &AppError{
		Code:    http.StatusForbidden,
		Message: ForbiddenMessage,
		Err:     err,
	}
}

func NotFound(err error) error {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: NotFoundMessage,
		Err:     err,
	}
}

func Conflict(err error) error {
	return &AppError{
		Code:    http.StatusConflict,
		Message: ConflictMessage,
		Err:     err,
	}
}

func TooManyRequests(err error) error {
	return &AppError{Code: http.StatusTooManyRequests, Message: TooManyRequestsMessage, Err: err}
}

func GatewayTimeout(err error) error {
	return &AppError{
		Code:    http.StatusGatewayTimeout,
		Message: GatewayTimeOutMessage,
		Err:     err,
	}
}
