package logging_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DoWithLogic/golang-clean-architecture/pkg/logging"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// newTestLogger returns a zerolog.Logger that writes JSON to the provided buffer.
func newTestLogger(buf *bytes.Buffer) *zerolog.Logger {
	l := zerolog.New(buf)
	return &l
}

// newEchoContext builds a minimal Echo context around an httptest.Request.
func newEchoContext(method, path string, body io.Reader, headers map[string]string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, path, body)
	req.Header.Set(echo.HeaderContentType, "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

// loggedFields unmarshals the first JSON log line from buf.
func loggedFields(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()
	var fields map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &fields); err != nil {
		t.Fatalf("failed to parse log output %q: %v", buf.String(), err)
	}
	return fields
}

func TestMiddleware_LogsRequest(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := newTestLogger(buf)

	c, _ := newEchoContext(http.MethodGet, "/api/test", nil, nil)

	mw := logging.Middleware(logging.WithLogger(logger))
	handler := mw(func(c echo.Context) error {
		return c.String(http.StatusOK, `{"ok":true}`)
	})

	if err := handler(c); err != nil {
		t.Fatalf("handler error: %v", err)
	}

	fields := loggedFields(t, buf)
	if fields["method"] != "GET" {
		t.Errorf("method = %v, want GET", fields["method"])
	}
	if fields["uri"] != "/api/test" {
		t.Errorf("uri = %v, want /api/test", fields["uri"])
	}
}

func TestMiddleware_SetsRequestIDHeader(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := newTestLogger(buf)

	c, rec := newEchoContext(http.MethodGet, "/", nil, nil)

	mw := logging.Middleware(logging.WithLogger(logger))
	handler := mw(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	handler(c)

	if rec.Header().Get("X-Request-ID") == "" {
		t.Error("expected X-Request-ID response header to be set")
	}
}

func TestMiddleware_PropagatesExistingRequestID(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := newTestLogger(buf)

	c, rec := newEchoContext(http.MethodGet, "/", nil, map[string]string{
		"X-Request-ID": "existing-id",
	})

	mw := logging.Middleware(logging.WithLogger(logger))
	handler := mw(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	handler(c)

	if got := rec.Header().Get("X-Request-ID"); got != "existing-id" {
		t.Errorf("X-Request-ID = %q, want %q", got, "existing-id")
	}
}

func TestMiddleware_StoresRequestIDInContext(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := newTestLogger(buf)

	c, _ := newEchoContext(http.MethodGet, "/", nil, map[string]string{
		"X-Request-ID": "ctx-id",
	})

	var capturedID any
	mw := logging.Middleware(logging.WithLogger(logger))
	handler := mw(func(c echo.Context) error {
		capturedID = c.Request().Context().Value(logging.RequestIDContextKey)
		return c.NoContent(http.StatusOK)
	})
	handler(c)

	if capturedID != "ctx-id" {
		t.Errorf("context request ID = %v, want ctx-id", capturedID)
	}
}

func TestMiddleware_SkipsIgnoredPath(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := newTestLogger(buf)

	ip, _ := logging.NewIgnoredPatterns([]string{`^/ping$`})

	c, _ := newEchoContext(http.MethodGet, "/ping", nil, nil)

	mw := logging.Middleware(logging.WithLogger(logger), logging.WithIgnoredPatterns(ip))
	handler := mw(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	handler(c)

	if buf.Len() > 0 {
		t.Errorf("expected no log output for ignored path, got: %s", buf.String())
	}
}

func TestMiddleware_LogsStatus(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := newTestLogger(buf)

	c, _ := newEchoContext(http.MethodGet, "/", nil, nil)

	mw := logging.Middleware(logging.WithLogger(logger))
	handler := mw(func(c echo.Context) error {
		return c.NoContent(http.StatusCreated)
	})
	handler(c)

	fields := loggedFields(t, buf)
	// JSON numbers unmarshal as float64.
	if status, ok := fields["status"].(float64); !ok || int(status) != http.StatusCreated {
		t.Errorf("status = %v, want %d", fields["status"], http.StatusCreated)
	}
}

func TestMiddleware_LogsRequestBody(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := newTestLogger(buf)

	body := `{"name":"dave"}`
	c, _ := newEchoContext(http.MethodPost, "/users", strings.NewReader(body), nil)

	mw := logging.Middleware(logging.WithLogger(logger))
	handler := mw(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	handler(c)

	fields := loggedFields(t, buf)
	reqBody, ok := fields["request_body"].(map[string]any)
	if !ok {
		t.Fatalf("request_body is not a map: %v", fields["request_body"])
	}
	if reqBody["name"] != "dave" {
		t.Errorf("request_body.name = %v, want dave", reqBody["name"])
	}
}

func TestMiddleware_MasksRequestBody(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := newTestLogger(buf)

	body := `{"username":"alice","password":"secret"}`
	c, _ := newEchoContext(http.MethodPost, "/login", strings.NewReader(body), nil)

	mw := logging.Middleware(logging.WithLogger(logger), logging.WithMaskedKeys("password"))
	handler := mw(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	handler(c)

	fields := loggedFields(t, buf)
	reqBody := fields["request_body"].(map[string]any)
	if reqBody["password"] != "*****" {
		t.Errorf("password should be masked in log, got %v", reqBody["password"])
	}
}

func TestMiddleware_LogsResponseBody(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := newTestLogger(buf)

	c, _ := newEchoContext(http.MethodGet, "/", nil, nil)

	mw := logging.Middleware(logging.WithLogger(logger))
	handler := mw(func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, "application/json")
		return c.JSON(http.StatusOK, map[string]string{"result": "ok"})
	})
	handler(c)

	fields := loggedFields(t, buf)
	respBody, ok := fields["response_body"].(map[string]any)
	if !ok {
		t.Fatalf("response_body is not a map: %v", fields["response_body"])
	}
	if respBody["result"] != "ok" {
		t.Errorf("response_body.result = %v, want ok", respBody["result"])
	}
}

func TestMiddleware_OmitsResponseBodyForNonJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := newTestLogger(buf)

	c, _ := newEchoContext(http.MethodGet, "/", nil, nil)

	mw := logging.Middleware(logging.WithLogger(logger))
	handler := mw(func(c echo.Context) error {
		return c.String(http.StatusOK, "plain text")
	})
	handler(c)

	fields := loggedFields(t, buf)
	if _, exists := fields["response_body"]; exists {
		t.Error("response_body should be absent for non-JSON responses")
	}
}

func TestMiddleware_LogsLatency(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := newTestLogger(buf)

	c, _ := newEchoContext(http.MethodGet, "/", nil, nil)

	mw := logging.Middleware(logging.WithLogger(logger))
	handler := mw(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	handler(c)

	fields := loggedFields(t, buf)
	if _, ok := fields["latency"]; !ok {
		t.Error("expected latency field in log")
	}
	if _, ok := fields["latency_human"]; !ok {
		t.Error("expected latency_human field in log")
	}
}

func TestMiddleware_PropagatesHandlerError(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := newTestLogger(buf)

	c, _ := newEchoContext(http.MethodGet, "/", nil, nil)

	want := echo.ErrInternalServerError
	mw := logging.Middleware(logging.WithLogger(logger))
	handler := mw(func(c echo.Context) error {
		return want
	})

	err := handler(c)
	if err != want {
		t.Errorf("handler error = %v, want %v", err, want)
	}
}
