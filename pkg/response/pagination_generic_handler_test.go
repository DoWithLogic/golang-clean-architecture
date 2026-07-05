package response

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

type testPaginationRequest struct {
	Name string `json:"name"`
}

func (r testPaginationRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

type testPaginationResponse struct {
	ID int `json:"id"`
}

func TestGenericPaginationHandler(t *testing.T) {
	t.Run("bind error", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		var request testPaginationRequest

		err := GenericPaginationHandler(
			ctx,
			&request,
			func(ctx context.Context, req *testPaginationRequest) (testPaginationResponse, error) {
				t.Fatal("usecase should not be called")
				return testPaginationResponse{}, nil
			},
		)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("validation error", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		var request testPaginationRequest

		err := GenericPaginationHandler(
			ctx,
			&request,
			func(ctx context.Context, req *testPaginationRequest) (testPaginationResponse, error) {
				t.Fatal("usecase should not be called")
				return testPaginationResponse{}, nil
			},
		)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("usecase error", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"john"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		var request testPaginationRequest

		err := GenericPaginationHandler(
			ctx,
			&request,
			func(ctx context.Context, req *testPaginationRequest) (testPaginationResponse, error) {
				return testPaginationResponse{}, BadRequest(errors.New("invalid"))
			},
		)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"john"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		var request testPaginationRequest

		err := GenericPaginationHandler(
			ctx,
			&request,
			func(ctx context.Context, req *testPaginationRequest) (testPaginationResponse, error) {
				return testPaginationResponse{
					ID: 123,
				}, nil
			},
		)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}

		var resp SuccessWithPaginationResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal response: %v", err)
		}

		if resp.Code != http.StatusOK {
			t.Fatalf("Code = %d", resp.Code)
		}

		if resp.Message != SuccessMessage {
			t.Fatalf("Message = %v", resp.Message)
		}

		data, ok := resp.Data.(map[string]any)
		if !ok {
			t.Fatalf("unexpected data type %T", resp.Data)
		}

		if data["id"].(float64) != 123 {
			t.Fatalf("unexpected response data: %+v", data)
		}
	})
}
