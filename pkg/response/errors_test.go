package response

import (
	"errors"
	"net/http"
	"testing"
)

func TestAppError_Error(t *testing.T) {
	err := errors.New("something went wrong")

	appErr := &AppError{
		Code:    http.StatusBadRequest,
		Message: BadRequestMessage,
		Err:     err,
	}

	if got := appErr.Error(); got != err.Error() {
		t.Errorf("Error() = %q, want %q", got, err.Error())
	}
}

func TestAppError_Unwrap(t *testing.T) {
	err := errors.New("wrapped error")

	appErr := &AppError{
		Err: err,
	}

	if got := appErr.Unwrap(); got != err {
		t.Errorf("Unwrap() = %v, want %v", got, err)
	}

	if !errors.Is(appErr, err) {
		t.Error("errors.Is() should identify wrapped error")
	}
}

func TestErrorConstructors(t *testing.T) {
	baseErr := errors.New("base error")

	tests := []struct {
		name    string
		fn      func(error) error
		code    int
		message ResponseMessage
	}{
		{
			name:    "BadRequest",
			fn:      BadRequest,
			code:    http.StatusBadRequest,
			message: BadRequestMessage,
		},
		{
			name:    "InternalServerError",
			fn:      InternalServerError,
			code:    http.StatusInternalServerError,
			message: InternalServerErrorMessage,
		},
		{
			name:    "Unauthorized",
			fn:      Unauthorized,
			code:    http.StatusUnauthorized,
			message: UnauthorizedMessage,
		},
		{
			name:    "Forbidden",
			fn:      Forbidden,
			code:    http.StatusForbidden,
			message: ForbiddenMessage,
		},
		{
			name:    "NotFound",
			fn:      NotFound,
			code:    http.StatusNotFound,
			message: NotFoundMessage,
		},
		{
			name:    "Conflict",
			fn:      Conflict,
			code:    http.StatusConflict,
			message: ConflictMessage,
		},
		{
			name:    "TooManyRequests",
			fn:      TooManyRequests,
			code:    http.StatusTooManyRequests,
			message: TooManyRequestsMessage,
		},
		{
			name:    "GatewayTimeout",
			fn:      GatewayTimeout,
			code:    http.StatusGatewayTimeout,
			message: GatewayTimeOutMessage,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn(baseErr)

			appErr, ok := err.(*AppError)
			if !ok {
				t.Fatalf("expected *AppError, got %T", err)
			}

			if appErr.Code != tt.code {
				t.Errorf("Code = %d, want %d", appErr.Code, tt.code)
			}

			if appErr.Message != tt.message {
				t.Errorf("Message = %+v, want %+v", appErr.Message, tt.message)
			}

			if appErr.Err != baseErr {
				t.Errorf("Err = %v, want %v", appErr.Err, baseErr)
			}

			if !errors.Is(appErr, baseErr) {
				t.Error("errors.Is() should identify wrapped error")
			}
		})
	}
}
