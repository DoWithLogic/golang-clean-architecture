package apperror

import (
	"errors"
	"net/http"
	"strings"
)

var (
	ErrEmailAlreadyExist = errors.New("email already exist")
	ErrInvalidUserType   = errors.New("invalid user_type")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrFailedGenerateJWT = errors.New("failed generate access token")
	ErrInvalidIsActive   = errors.New("invalid is_active")
	ErrStatusValue       = errors.New("status should be 0 or 1")

	ErrFailedGetTokenInformation = errors.New("failed to get token information")
)

type AppError struct {
	Code    int
	Err     error
	Message string
}

func Equals(err error, expectedErr error) bool {
	return strings.EqualFold(err.Error(), expectedErr.Error())
}

func (h AppError) Error() string {
	return h.Err.Error()
}

func BadRequest(err error) error {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: "bad_request",
		Err:     err,
	}
}

func InternalServerError(err error) error {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: "internal_server_error",
		Err:     err,
	}
}

func Unauthorized(err error) error {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: "unauthorized",
		Err:     err,
	}
}

func Forbidden(err error) error {
	return &AppError{
		Code:    http.StatusForbidden,
		Message: "forbidden",
		Err:     err,
	}
}

func NotFound(err error) error {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: "not_found",
		Err:     err,
	}
}

func Conflict(err error) error {
	return &AppError{
		Code:    http.StatusConflict,
		Message: "Conflict",
		Err:     err,
	}
}

func GatewayTimeout(err error) error {
	return &AppError{
		Code:    http.StatusGatewayTimeout,
		Message: "gateway_timeout",
		Err:     err,
	}
}
