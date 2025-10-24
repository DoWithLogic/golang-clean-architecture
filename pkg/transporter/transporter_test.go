package transporter_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DoWithLogic/golang-clean-architecture/pkg/constants"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/transporter"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestDoHttpRequestError(t *testing.T) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	appHttp := transporter.NewClient(&logger)

	// Create a test server that returns an error response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a 404 Not Found error
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found"}`))
	}))
	defer server.Close()

	req := transporter.Request{
		Method:   "GET",
		Endpoint: server.URL + "/nonexistent", // Use the test server URL
		Headers:  map[string]string{"Content-Type": constants.MIMEApplicationJSON},
	}

	var res struct {
		Error string `json:"error"`
	}

	// Call the DoHttpRequest method
	err := appHttp.DoHttpRequest(context.Background(), req, &res)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status code: 404")
}

func TestDoHttpRequestWithJSONBody(t *testing.T) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	appHttp := transporter.NewClient(&logger)

	// Create a test server that returns a mocked response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for the Content-Type header
		if r.Header.Get("Content-Type") != constants.MIMEApplicationJSON {
			http.Error(w, "Invalid content type", http.StatusBadRequest)
			return
		}

		var reqBody map[string]string
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Respond with the expected JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}))
	defer server.Close()

	req := transporter.Request{
		Method:   "POST",
		Endpoint: server.URL + "/test",
		Headers:  map[string]string{"Content-Type": constants.MIMEApplicationJSON},
		Body:     map[string]string{"key": "value"},
	}

	var res struct {
		Success bool `json:"success"`
	}

	// Call the DoHttpRequest method
	err := appHttp.DoHttpRequest(context.Background(), req, &res)

	// Assertions
	assert.NoError(t, err)
	assert.True(t, res.Success)
}

func TestDoHttpRequestWithFormFile(t *testing.T) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	appHttp := transporter.NewClient(&logger)

	// Create a test server that handles form file uploads
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse the multipart form
		err := r.ParseMultipartForm(10 << 20) // Limit your max memory
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// Check if the file exists
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "File not found", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Respond with a success message
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}))
	defer server.Close()

	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up

	// Write some test data to the file
	if _, err := tempFile.WriteString("This is a test file."); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Close the file so it can be read later
	defer tempFile.Close()

	// Prepare the request with the file
	req := transporter.Request{
		Method:   http.MethodPost,
		Endpoint: server.URL + "/upload",
		Headers:  map[string]string{"Content-Type": constants.MIMEMultipartForm},
		Files: map[string]transporter.File{
			"file": {
				FileName: "fileku.text",
				File:     tempFile,
			},
		},
	}

	var res struct {
		Success bool `json:"success"`
	}

	// Call the DoHttpRequest method
	err = appHttp.DoHttpRequest(context.Background(), req, &res)

	// Assertions
	assert.NoError(t, err)
	assert.True(t, res.Success)
}
