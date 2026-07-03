package logging_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DoWithLogic/golang-clean-architecture/pkg/logging"
)

func TestNewCustomResponseWriter_DefaultStatus(t *testing.T) {
	rec := httptest.NewRecorder()
	crw := logging.New(rec)

	if crw.Status() != http.StatusOK {
		t.Errorf("default status = %d, want %d", crw.Status(), http.StatusOK)
	}
}

func TestNewCustomResponseWriter_InitialSizeZero(t *testing.T) {
	rec := httptest.NewRecorder()
	crw := logging.New(rec)

	if crw.Size() != 0 {
		t.Errorf("initial size = %d, want 0", crw.Size())
	}
}

func TestNewCustomResponseWriter_InitialBodyEmpty(t *testing.T) {
	rec := httptest.NewRecorder()
	crw := logging.New(rec)

	if len(crw.Body()) != 0 {
		t.Errorf("initial body should be empty, got %q", crw.Body())
	}
}

func TestCustomResponseWriter_Write_BuffersBody(t *testing.T) {
	rec := httptest.NewRecorder()
	crw := logging.New(rec)

	payload := []byte(`{"message":"hello"}`)
	n, err := crw.Write(payload)

	if err != nil {
		t.Fatalf("Write() error: %v", err)
	}
	if n != len(payload) {
		t.Errorf("Write() n = %d, want %d", n, len(payload))
	}
	if string(crw.Body()) != string(payload) {
		t.Errorf("Body() = %q, want %q", crw.Body(), payload)
	}
}

func TestCustomResponseWriter_Write_ForwardsToUnderlying(t *testing.T) {
	rec := httptest.NewRecorder()
	crw := logging.New(rec)

	payload := []byte("forwarded")
	crw.Write(payload)

	if rec.Body.String() != "forwarded" {
		t.Errorf("underlying writer got %q, want %q", rec.Body.String(), "forwarded")
	}
}

func TestCustomResponseWriter_Write_AccumulatesSize(t *testing.T) {
	rec := httptest.NewRecorder()
	crw := logging.New(rec)

	crw.Write([]byte("abc"))
	crw.Write([]byte("de"))

	if crw.Size() != 5 {
		t.Errorf("Size() = %d, want 5", crw.Size())
	}
}

func TestCustomResponseWriter_Write_AccumulatesBody(t *testing.T) {
	rec := httptest.NewRecorder()
	crw := logging.New(rec)

	crw.Write([]byte("foo"))
	crw.Write([]byte("bar"))

	if string(crw.Body()) != "foobar" {
		t.Errorf("Body() = %q, want %q", crw.Body(), "foobar")
	}
}

func TestCustomResponseWriter_WriteHeader_CapturesStatus(t *testing.T) {
	rec := httptest.NewRecorder()
	crw := logging.New(rec)

	crw.WriteHeader(http.StatusCreated)

	if crw.Status() != http.StatusCreated {
		t.Errorf("Status() = %d, want %d", crw.Status(), http.StatusCreated)
	}
}

func TestCustomResponseWriter_WriteHeader_ForwardsToUnderlying(t *testing.T) {
	rec := httptest.NewRecorder()
	crw := logging.New(rec)

	crw.WriteHeader(http.StatusNotFound)

	if rec.Code != http.StatusNotFound {
		t.Errorf("underlying recorder code = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestCustomResponseWriter_WriteHeader_MultipleStatuses_LastWins(t *testing.T) {
	rec := httptest.NewRecorder()
	crw := logging.New(rec)

	// In real HTTP only the first WriteHeader takes effect on the wire,
	// but our crw.status field always reflects the last call.
	crw.WriteHeader(http.StatusOK)
	crw.WriteHeader(http.StatusInternalServerError)

	if crw.Status() != http.StatusInternalServerError {
		t.Errorf("Status() = %d, want %d", crw.Status(), http.StatusInternalServerError)
	}
}

func TestCustomResponseWriter_Hijack_NotSupported(t *testing.T) {
	// httptest.ResponseRecorder does NOT implement http.Hijacker.
	rec := httptest.NewRecorder()
	crw := logging.New(rec)

	_, _, err := crw.Hijack()
	if err == nil {
		t.Error("expected error when underlying writer does not support Hijack")
	}
}
