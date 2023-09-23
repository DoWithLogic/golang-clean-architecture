package dtos

import (
	"errors"
	"net/http"
)

type (
	Message map[string]string

	Response struct {
		Status  int            `json:"status"`
		Message Message        `json:"message"`
		Errors  []CaptureError `json:"errors,omitempty"`
		Data    interface{}    `json:"data,omitempty"`
		Meta    interface{}    `json:"meta,omitempty"`
		Header  http.Header    `json:"header,omitempty"`
		Body    interface{}    `json:"body,omitempty"`
	}

	CaptureError struct {
		Details string `json:"details"`
		Message string `json:"message"`
	}
)

var (
	Text = http.StatusText

	MsgSuccess = map[string]string{"en": "Success", "id": "Sukses"}
	MsgFailed  = map[string]string{"en": "Failed", "id": "Gagal"}
)

func UnwrapFirstError(err error) string { return UnwrapAll(err).Error() }

// UnwrapAll will unwrap the underlying error until we get the first wrapped error.
func UnwrapAll(err error) error {
	for err != nil && errors.Unwrap(err) != nil {
		err = errors.Unwrap(err)
	}

	return err
}

func NewResponse(statusCode int, message Message, data interface{}) Response {
	return Response{
		Status:  statusCode,
		Message: MsgSuccess,
		Data:    data,
	}
}

func NewResponseError(statusCode int, messageStatus Message, message, details string) Response {
	return Response{
		Status:  statusCode,
		Message: messageStatus,
		Errors: []CaptureError{
			{
				Message: message,
				Details: details,
			},
		},
	}
}
