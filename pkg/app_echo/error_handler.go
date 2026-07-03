package app_echo

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/DoWithLogic/golang-clean-architecture/pkg/response"
	"github.com/golang-jwt/jwt"
	"github.com/invopop/validation"
	"github.com/labstack/echo/v4"
	"github.com/valyala/fasthttp"
)

var (
	ErrAuthenticationRequired = errors.New("authentication required")
	ErrResourceNotFound       = errors.New("resource not found")
	ErrMalformedJSON          = errors.New("malformed JSON request body")
	ErrValidationFailed       = errors.New("request validation failed")
	ErrInvalidMultipart       = errors.New("invalid multipart/form-data request")
	ErrInvalidPathParameter   = errors.New("invalid path parameter: expected number")
	ErrMissingOrMalformatJWT  = errors.New("Missing or malformed JWT")
)

// errorHandler handles application errors and returns the appropriate HTTP response.
func errorHandler(err error, c echo.Context) {
	var jwtErr *jwt.ValidationError
	if errors.As(err, &jwtErr) || errors.Is(err, ErrMissingOrMalformatJWT) {
		response.ErrorBuilder(response.Unauthorized(ErrAuthenticationRequired)).Send(c)
		return
	}

	var echoErr *echo.HTTPError
	if errors.As(err, &echoErr) {
		switch echoErr.Code {
		case http.StatusNotFound:
			response.ErrorBuilder(response.NotFound(ErrResourceNotFound)).Send(c)
		default:
			response.ErrorBuilder(err).Send(c)
		}

		return
	}

	var numErr *strconv.NumError
	if errors.As(err, &numErr) {
		response.ErrorBuilder(response.BadRequest(ErrInvalidPathParameter)).Send(c)
		return
	}

	var appErr *response.AppError
	if errors.As(err, &appErr) {
		response.ErrorBuilder(appErr).Send(c)
		return
	}

	var validationErr validation.Errors
	if errors.As(err, &validationErr) {
		response.ErrorBuilder(response.BadRequest(ErrValidationFailed)).Send(c)
		return
	}

	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		response.ErrorBuilder(response.BadRequest(ErrMalformedJSON)).Send(c)
		return
	}

	var unmarshalErr *json.UnmarshalTypeError
	if errors.As(err, &unmarshalErr) {
		field := unmarshalErr.Field
		if field == "" {
			field = "<unknown>"
		}

		response.ErrorBuilder(
			response.BadRequest(
				fmt.Errorf(
					`invalid value for field "%s": expected %s`,
					field,
					expectedType(unmarshalErr.Type),
				),
			),
		).Send(c)

		return
	}

	var timeErr *time.ParseError
	if errors.As(err, &timeErr) {
		value := timeErr.Value
		if value == "" {
			value = "empty string"
		}

		response.ErrorBuilder(
			response.BadRequest(
				fmt.Errorf(
					`invalid datetime "%s": expected ISO 8601 (RFC 3339), e.g. 2026-07-04T15:30:00Z`,
					value,
				),
			),
		).Send(c)

		return
	}

	if errors.Is(err, fasthttp.ErrNoMultipartForm) {
		response.ErrorBuilder(response.BadRequest(ErrInvalidMultipart)).Send(c)
		return
	}

	var netErr *net.OpError
	if errors.As(err, &netErr) {
		log.Fatalf("unable to establish TCP connection to %s", netErr.Addr)
	}

	response.ErrorBuilder(err).Send(c)
}

// expectedType returns a human-readable type name.
func expectedType(t reflect.Type) string {
	if t == reflect.TypeFor[time.Time]() {
		return "ISO 8601 datetime (RFC 3339)"
	}

	switch t.Kind() {
	case reflect.Bool:
		return "boolean"

	case reflect.String:
		return "string"

	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64:
		return "number"

	case reflect.Slice, reflect.Array:
		return "array"

	case reflect.Map, reflect.Struct:
		return "object"

	default:
		return t.String()
	}
}
