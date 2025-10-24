package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"strconv"
	"time"

	"github.com/DoWithLogic/golang-clean-architecture/pkg/errs"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/response"
	"github.com/golang-jwt/jwt"
	"github.com/invopop/validation"
	"github.com/labstack/echo/v4"

	"github.com/valyala/fasthttp"
)

// errorHandler handles HTTP errors and sends a custom response based on the error type.
// Parameters:
//   - err: The error that occurred.
//   - c: The Echo context.
func errorHandler(err error, c echo.Context) {
	// Unauthorized Error
	var jwtErr *jwt.ValidationError
	if errors.As(err, &jwtErr) || err.Error() == "Missing or malformed JWT" {
		response.ErrorBuilder(errs.Unauthorized(err)).Send(c)

		return
	}

	var echoErr *echo.HTTPError
	if errors.As(err, &echoErr) {
		report, ok := err.(*echo.HTTPError)

		if !ok {
			report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		switch report.Code {
		case http.StatusNotFound:
			response.ErrorBuilder(errs.NotFound(errors.New("route not found"))).Send(c)

			return
		default:
			response.ErrorBuilder(err).Send(c)

			return
		}
	}

	// Path Parse Error
	var numErr *strconv.NumError
	if errors.As(err, &numErr) {
		response.ErrorBuilder(errs.BadRequest(errors.New("malformed_body"))).Send(c)

		return
	}

	// handle HTTP Error
	var appErr *errs.AppError
	if errors.As(err, &appErr) {
		response.ErrorBuilder(err).Send(c)

		return
	}

	var validatorError validation.Errors
	if errors.As(err, &validatorError) {
		response.ErrorBuilder(errs.BadRequest(errors.New("validation_error"))).Send(c)
		return
	}

	// JSON Format Error
	var jsonSyntaxErr *json.SyntaxError
	if errors.As(err, &jsonSyntaxErr) {
		response.ErrorBuilder(errs.BadRequest(errors.New("malformed_body"))).Send(c)

		return
	}

	// Unmarshal Error
	var unmarshalErr *json.UnmarshalTypeError
	if errors.As(err, &unmarshalErr) {
		var translatedType string
		switch unmarshalErr.Type.Name() {
		// REGEX *int*
		case "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64", "float32", "float64":
			translatedType = "number"
		case "Time":
			translatedType = "date time"
		case "string":
			translatedType = "string"
		}

		response.ErrorBuilder(errs.BadRequest(fmt.Errorf("the field must be a valid %s", translatedType))).Send(c)
		return
	}

	//time parse error
	var timeParseErr *time.ParseError
	if errors.As(err, &timeParseErr) {
		v := timeParseErr.Value
		if v == "" {
			v = "empty string (``)"
		}

		response.ErrorBuilder(errs.BadRequest(fmt.Errorf("invalid time format on %s", v))).Send(c)

		return
	}

	// Multipart Error
	if errors.Is(err, fasthttp.ErrNoMultipartForm) {
		response.ErrorBuilder(errs.BadRequest(errors.New("invalid multipart content-type"))).Send(c)

		return
	}

	//TCP connection error
	var tcpErr *net.OpError
	if errors.As(err, &tcpErr) {
		log.Fatalf("unable to get tcp connection from %s, shutting down...", tcpErr.Addr.String())
	}

	response.ErrorBuilder(err).Send(c)
}
