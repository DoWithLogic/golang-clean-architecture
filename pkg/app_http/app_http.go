package app_http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"reflect"
	"time"

	"github.com/DoWithLogic/golang-clean-architecture/pkg/constant"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/observability/instrumentation"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
)

type (
	File struct {
		FileName string
		File     io.Reader
	}

	Request struct {
		Method   string
		Endpoint string
		Headers  map[string]string
		Body     any
		Files    map[string]File
	}
)

// validateBody checks if the provided body is a non-pointer type
func (r Request) validateBody() error {
	if r.Body == nil {
		return nil
	}
	v := reflect.ValueOf(r.Body)
	// Check if the body is a pointer
	if v.Kind() == reflect.Ptr {
		return errors.New("body must be a non-pointer type")
	}
	return nil
}

type AppHttp struct {
	client *fasthttp.Client
	log    *zerolog.Logger
}

// NewClient creates a new fasthttp client with default settings
func NewClient(log *zerolog.Logger) *AppHttp {
	return &AppHttp{
		client: &fasthttp.Client{
			// Maximum number of connections allowed per host. This controls the number of keep-alive connections.
			MaxConnsPerHost: 50,
			// The function used to establish network connections. The default is sufficient for most cases.
			Dial: fasthttp.Dial,
			// The maximum time a connection can remain idle before being closed.
			MaxIdleConnDuration: 30 * time.Second,
			// Maximum time allowed for reading a response from the server.
			ReadTimeout: 10 * time.Second,
			// Maximum time allowed for writing a request to the server.
			WriteTimeout: 10 * time.Second,
		},
		log: log,
	}
}
func (c *AppHttp) DoHttpRequest(ctx context.Context, req Request, res any) error {
	ctx, span := instrumentation.NewTraceSpan(ctx, "DoHttpRequest")
	defer span.End()

	if err := req.validateBody(); err != nil {
		return err
	}

	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(request)
		fasthttp.ReleaseResponse(response)
	}()

	request.Header.SetMethod(req.Method)
	request.SetRequestURI(req.Endpoint)

	// Set request headers
	for key, value := range req.Headers {
		request.Header.Set(key, value)
	}
	// If there are files to upload, create multipart form data
	if req.Files != nil {
		var buffer bytes.Buffer
		writer := multipart.NewWriter(&buffer)

		// Add files to the form
		for key, file := range req.Files {
			part, err := writer.CreateFormFile(key, file.FileName) // Adjust filename as needed
			if err != nil {
				return errors.Wrap(err, "failed to create form file")
			}

			if _, err := io.Copy(part, file.File); err != nil {
				return errors.Wrap(err, "failed to copy file to form")
			}
		}

		// Close the writer to finalize the form data
		if err := writer.Close(); err != nil {
			return errors.Wrap(err, "failed to close writer")
		}

		request.SetBody(buffer.Bytes())
		request.Header.Set("Content-Type", writer.FormDataContentType())
	} else if req.Body != nil {
		// If there is a body, marshal it to JSON
		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			c.log.Err(err).Ctx(ctx).Msg("[DoHttpRequest]Marshal")
			return errors.Wrap(err, "failed to marshal request body")
		}

		request.SetBody(jsonBody)
		request.Header.Set("Content-Type", constant.MIMEApplicationJSON)
	}

	start := time.Now()
	c.log.Info().Ctx(ctx).
		Str("method", req.Method).
		Str("endpoint", req.Endpoint).
		Interface("headers", req.Headers).
		Interface("body", string(request.Body())).
		Msg("[DoHttpRequest]Sending request")

	// Perform the request
	if err := c.client.Do(request, response); err != nil {
		c.log.Err(err).Ctx(ctx).Msg("[DoHttpRequest]client.Do")
		return errors.Wrap(err, "failed to execute HTTP request")
	}

	c.log.Info().Ctx(ctx).
		Int("status_code", response.StatusCode()).
		Dur("duration", time.Since(start)).
		RawJSON("response", response.Body()).
		Msg("[DoHttpRequest]Received response")

	// Check response status code
	if response.StatusCode() != fasthttp.StatusOK {
		c.log.Error().Ctx(ctx).
			Int("status_code", response.StatusCode()).
			Dur("duration", time.Since(start)).
			Msg("[DoHttpRequest] Unexpected status code")

		return fmt.Errorf("unexpected status code: %v", response.StatusCode())
	}

	// Decode response if a response struct is provided
	if response != nil {
		if err := json.Unmarshal(response.Body(), res); err != nil {
			c.log.Err(err).Ctx(ctx).Msg("[DoHttpRequest]json.Unmarshal")
			return errors.Wrap(err, "failed to decode response")
		}
	}
	return nil
}
