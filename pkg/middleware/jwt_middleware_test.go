package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DoWithLogic/golang-clean-architecture/pkg/jwt"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/middleware"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/redis"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/types"
	"github.com/alicebob/miniredis"
	"github.com/labstack/echo/v4"
)

func newMiddleware(t *testing.T) (*middleware.Middleware, string) {
	t.Helper()

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer mr.Close()

	cfg := redis.RedisConfig{Addr: mr.Addr()}

	jwtFactory := jwt.NewJWTFactory(
		jwt.JWTConfig{
			Key:             "secret-key",
			ExpiredInSecond: 3600,
		},
		redis.NewRedisManager(redis.NewRedisClient(t.Context(), cfg)),
	)

	token, err := jwtFactory.CreateJWT(&jwt.JWTClaims{
		Data: &jwt.Data{
			ID:           1,
			ContactType:  types.CONTACT_TYPE_EMAIL,
			ContactValue: "john@example.com",
		},
	})
	if err != nil {
		t.Fatalf("CreateJWT() error = %v", err)
	}

	return middleware.New(jwtFactory), token
}

func TestJWTMiddleware_MissingAuthorizationHeader(t *testing.T) {
	m, _ := newMiddleware(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	nextCalled := false

	handler := m.JWTMiddleware()(func(c echo.Context) error {
		nextCalled = true
		return c.NoContent(http.StatusOK)
	})

	if err := handler(ctx); err != nil {
		t.Fatalf("handler() error = %v", err)
	}

	if nextCalled {
		t.Fatal("next handler should not be called")
	}

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestJWTMiddleware_InvalidToken(t *testing.T) {
	m, _ := newMiddleware(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(types.AuthorizationHeaderKey.String(), "Bearer invalid-token")
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	nextCalled := false

	handler := m.JWTMiddleware()(func(c echo.Context) error {
		nextCalled = true
		return c.NoContent(http.StatusOK)
	})

	if err := handler(ctx); err != nil {
		t.Fatalf("handler() error = %v", err)
	}

	if nextCalled {
		t.Fatal("next handler should not be called")
	}

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestJWTMiddleware_Success(t *testing.T) {
	m, token := newMiddleware(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(types.AuthorizationHeaderKey.String(), "Bearer "+token)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	nextCalled := false

	handler := m.JWTMiddleware()(func(c echo.Context) error {
		nextCalled = true

		v := c.Get(types.CredentialDataContextKey.String())
		if v == nil {
			t.Fatal("claims not found in context")
		}

		claims, ok := v.(*jwt.JWTClaims)
		if !ok {
			t.Fatalf("unexpected type %T", v)
		}

		if claims.Data.ID != 1 {
			t.Fatalf("ID = %d, want 1", claims.Data.ID)
		}

		if claims.Data.ContactValue != "john@example.com" {
			t.Fatalf("ContactValue = %s", claims.Data.ContactValue)
		}

		return c.NoContent(http.StatusOK)
	})

	if err := handler(ctx); err != nil {
		t.Fatalf("handler() error = %v", err)
	}

	if !nextCalled {
		t.Fatal("next handler should be called")
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}
