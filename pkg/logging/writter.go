// Package logging provides a reusable Echo middleware for structured HTTP
// request/response logging via zerolog. It supports sensitive-field masking,
// path exclusion, and X-Request-ID propagation.
//
// # Quick Start
//
//	import (
//	    logging "github.com/golang-clean-architecture/pkg/logging"
//	)
//
//	e := echo.New()
//	e.Use(logging.Middleware(
//	    logging.WithLogger(&logger),
//	    logging.WithMaskedKeys("password", "token", "secret"),
//	))
package logging

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
)

// CustomResponseWriter wraps http.ResponseWriter to intercept and buffer
// the response body and status code without disrupting the normal write flow.
type CustomResponseWriter struct {
	http.ResponseWriter
	buf          *bytes.Buffer
	status       int
	bytesWritten int
}

// New creates a new CustomResponseWriter wrapping the provided http.ResponseWriter.
// The default status code is 200 OK.
func New(w http.ResponseWriter) *CustomResponseWriter {
	return &CustomResponseWriter{
		ResponseWriter: w,
		buf:            new(bytes.Buffer),
		status:         http.StatusOK,
	}
}

// Write writes data to both the internal buffer and the underlying ResponseWriter.
// The buffer is used later for logging the response body.
func (crw *CustomResponseWriter) Write(b []byte) (int, error) {
	n, err := crw.buf.Write(b)
	if err != nil {
		return n, err
	}
	crw.bytesWritten += n
	return crw.ResponseWriter.Write(b)
}

// WriteHeader captures the HTTP status code and delegates to the underlying writer.
func (crw *CustomResponseWriter) WriteHeader(statusCode int) {
	crw.status = statusCode
	crw.ResponseWriter.WriteHeader(statusCode)
}

// Hijack implements http.Hijacker to support WebSocket upgrades.
// Returns an error if the underlying ResponseWriter does not support hijacking.
func (crw *CustomResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := crw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("ResponseWriter does not implement http.Hijacker")
	}
	return hijacker.Hijack()
}

// Body returns the buffered response body bytes.
func (crw *CustomResponseWriter) Body() []byte { return crw.buf.Bytes() }

// Status returns the HTTP status code written to the response.
func (crw *CustomResponseWriter) Status() int { return crw.status }

// Size returns the total number of bytes written to the response body.
func (crw *CustomResponseWriter) Size() int { return crw.bytesWritten }
