package errs_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/DoWithLogic/golang-clean-architecture/pkg/errs"
)

func TestEquals(t *testing.T) {
	tests := []struct {
		name     string
		err1     error
		err2     error
		expected bool
	}{
		{"equal errors (same case)", errors.New("error message"), errors.New("error message"), true},
		{"equal errors (different case)", errors.New("Error Message"), errors.New("error message"), true},
		{"different errors", errors.New("error A"), errors.New("error B"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := errs.Equals(tt.err1, tt.err2)
			if got != tt.expected {
				t.Errorf("Equals() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAppErrorMethods(t *testing.T) {
	baseErr := errors.New("base error")
	appErr := &errs.AppError{
		Code:    http.StatusBadRequest,
		Err:     baseErr,
		Message: "bad_request",
	}

	if appErr.Error() != "base error" {
		t.Errorf("Error() = %v, want %v", appErr.Error(), "base error")
	}

	if appErr.Unwrap() != baseErr {
		t.Errorf("Unwrap() = %v, want %v", appErr.Unwrap(), baseErr)
	}
}

func TestErrorHelpers(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(error) error
		err      error
		expected int
		message  string
	}{
		{"BadRequest", errs.BadRequest, errors.New("bad req"), http.StatusBadRequest, "bad_request"},
		{"InternalServerError", errs.InternalServerError, errors.New("internal"), http.StatusInternalServerError, "internal_server_error"},
		{"Unauthorized", errs.Unauthorized, errors.New("unauth"), http.StatusUnauthorized, "unauthorized"},
		{"Forbidden", errs.Forbidden, errors.New("forbid"), http.StatusForbidden, "forbidden"},
		{"NotFound", errs.NotFound, errors.New("not found"), http.StatusNotFound, "not_found"},
		{"Conflict", errs.Conflict, errors.New("conflict"), http.StatusConflict, "Conflict"},
		{"GatewayTimeout", errs.GatewayTimeout, errors.New("timeout"), http.StatusGatewayTimeout, "gateway_timeout"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn(tt.err)
			appErr, ok := got.(*errs.AppError)
			if !ok {
				t.Fatalf("Expected *AppError, got %T", got)
			}
			if appErr.Code != tt.expected {
				t.Errorf("Code = %v, want %v", appErr.Code, tt.expected)
			}
			if appErr.Message != tt.message {
				t.Errorf("Message = %v, want %v", appErr.Message, tt.message)
			}
			if appErr.Err != tt.err {
				t.Errorf("Err = %v, want %v", appErr.Err, tt.err)
			}
		})
	}
}
