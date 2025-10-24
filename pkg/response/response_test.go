package response_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DoWithLogic/golang-clean-architecture/pkg/errs"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/response"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// --- BasicBuilder ---
func TestBasicBuilder(t *testing.T) {
	input := response.BasicResponse{Code: 200, Message: "ok", Data: "data"}
	got := response.BasicBuilder(input)
	assert.Equal(t, input, got)
}

// --- ErrorBuilder ---
func TestErrorBuilderWithAppError(t *testing.T) {
	err := errs.BadRequest(errors.New("bad input"))
	got := response.ErrorBuilder(err)

	appErr := err.(*errs.AppError)
	assert.Equal(t, appErr.Code, got.Code)
	assert.Equal(t, appErr.Message, got.Message)
	assert.Equal(t, appErr.Error(), got.Error)
}

func TestErrorBuilderWithGenericError(t *testing.T) {
	err := errors.New("something failed")
	got := response.ErrorBuilder(err)
	assert.Equal(t, http.StatusInternalServerError, got.Code)
	assert.Equal(t, response.INTERNAL_SERVER_ERROR, got.Message)
	assert.Equal(t, "something failed", got.Error)
}

func TestErrorBuilderWithNilError(t *testing.T) {
	got := response.ErrorBuilder(nil)
	assert.Equal(t, http.StatusInternalServerError, got.Code)
	assert.Equal(t, response.INTERNAL_SERVER_ERROR, got.Message)
	assert.Equal(t, response.INTERNAL_SERVER_ERROR, got.Error)
}

// --- Send Methods ---
func TestBasicResponseSend(t *testing.T) {
	e := echo.New()
	rec := httptest.NewRecorder()
	ctx := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), rec)

	resp := response.BasicResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    "ok",
	}
	err := resp.Send(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestFailedResponseSend(t *testing.T) {
	e := echo.New()
	rec := httptest.NewRecorder()
	ctx := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), rec)

	resp := response.FailedResponse{
		Code:    http.StatusInternalServerError,
		Message: "error",
		Error:   "something failed",
	}
	err := resp.Send(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestSuccessResponseSend(t *testing.T) {
	e := echo.New()
	rec := httptest.NewRecorder()
	ctx := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), rec)

	resp := response.SuccessBuilder("payload", map[string]string{"page": "1"})
	err := resp.Send(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// --- SuccessBuilder ---
func TestSuccessBuilderNoMeta(t *testing.T) {
	got := response.SuccessBuilder("payload")
	assert.Equal(t, http.StatusOK, got.Code)
	assert.Equal(t, "success", got.Message)
	assert.Equal(t, "payload", got.Data)
	assert.Nil(t, got.Meta.Meta)
}

func TestSuccessBuilderWithMeta(t *testing.T) {
	meta := map[string]string{"page": "1"}
	got := response.SuccessBuilder("payload", meta)
	assert.Equal(t, meta, got.Meta.Meta)
}
