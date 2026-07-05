package response

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestErrorBuilder(t *testing.T) {
	t.Run("returns app error response", func(t *testing.T) {
		baseErr := errors.New("invalid request")

		resp := ErrorBuilder(BadRequest(baseErr))

		if resp.Code != http.StatusBadRequest {
			t.Fatalf("Code = %d, want %d", resp.Code, http.StatusBadRequest)
		}

		if resp.Message != BadRequestMessage {
			t.Fatalf("Message = %v, want %v", resp.Message, BadRequestMessage)
		}

		if resp.Error != baseErr.Error() {
			t.Fatalf("Error = %q, want %q", resp.Error, baseErr.Error())
		}
	})

	t.Run("returns internal server error for normal error", func(t *testing.T) {
		baseErr := errors.New("something went wrong")

		resp := ErrorBuilder(baseErr)

		if resp.Code != http.StatusInternalServerError {
			t.Fatalf("Code = %d, want %d", resp.Code, http.StatusInternalServerError)
		}

		if resp.Message != InternalServerErrorMessage {
			t.Fatalf("Message = %v, want %v", resp.Message, InternalServerErrorMessage)
		}

		if resp.Error != baseErr.Error() {
			t.Fatalf("Error = %q, want %q", resp.Error, baseErr.Error())
		}
	})

	t.Run("returns default internal server error for nil error", func(t *testing.T) {
		resp := ErrorBuilder(nil)

		if resp.Code != http.StatusInternalServerError {
			t.Fatalf("Code = %d, want %d", resp.Code, http.StatusInternalServerError)
		}

		if resp.Message != InternalServerErrorMessage {
			t.Fatalf("Message = %v, want %v", resp.Message, InternalServerErrorMessage)
		}

		if resp.Error != InternalServerErrorMessage.String() {
			t.Fatalf("Error = %q, want %q", resp.Error, InternalServerErrorMessage.String())
		}
	})
}

func TestErrorResponse_Send(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	resp := ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: BadRequestMessage,
		Error:   "invalid request",
	}

	if err := resp.Send(ctx); err != nil {
		t.Fatalf("Send() returned error: %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}

	expected := `{"code":400,"message":"bad_request","error":"invalid request"}`

	if rec.Body.String() != expected+"\n" {
		t.Fatalf("body = %s, want %s", rec.Body.String(), expected)
	}
}
