package response

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestSuccessBuilder(t *testing.T) {
	t.Run("without meta", func(t *testing.T) {
		data := map[string]any{
			"id":   1,
			"name": "john",
		}

		resp := SuccessBuilder(data)

		if resp.Code != http.StatusOK {
			t.Fatalf("Code = %d, want %d", resp.Code, http.StatusOK)
		}

		if resp.Message != SuccessMessage {
			t.Fatalf("Message = %v, want %v", resp.Message, SuccessMessage)
		}

		if !reflect.DeepEqual(resp.Data, data) {
			t.Fatalf("Data = %#v, want %#v", resp.Data, data)
		}

		if resp.Meta != nil {
			t.Fatalf("Meta = %#v, want nil", resp.Meta)
		}
	})

	t.Run("with meta", func(t *testing.T) {
		data := []string{"foo", "bar"}
		meta := map[string]any{
			"page":  1,
			"limit": 10,
			"total": 20,
		}

		resp := SuccessBuilder(data, meta)

		if resp.Code != http.StatusOK {
			t.Fatalf("Code = %d, want %d", resp.Code, http.StatusOK)
		}

		if resp.Message != SuccessMessage {
			t.Fatalf("Message = %v, want %v", resp.Message, SuccessMessage)
		}

		if !reflect.DeepEqual(resp.Data, data) {
			t.Fatalf("Data = %#v, want %#v", resp.Data, data)
		}

		if !reflect.DeepEqual(resp.Meta, meta) {
			t.Fatalf("Meta = %#v, want %#v", resp.Meta, meta)
		}
	})
}

func TestSuccessResponse_Send(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	resp := SuccessResponse{
		SuccessDefault: SuccessDefault{
			Code:    http.StatusOK,
			Message: SuccessMessage,
		},
		Data: map[string]any{
			"id": 1,
		},
	}

	if err := resp.Send(ctx); err != nil {
		t.Fatalf("Send() returned error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	expected := `{"code":200,"message":"success","data":{"id":1}}`

	if rec.Body.String() != expected+"\n" {
		t.Fatalf("body = %s, want %s", rec.Body.String(), expected)
	}
}
