package data_test

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"ohcl/controllers/data"
	"os"
	"path/filepath"
	"testing"
)

// TestPostEndpoint will test post endpoint.
func TestPostEndpoint(t *testing.T) {
	// Set up the test server.
	r := gin.Default()
	r.POST("/data", data.Post)

	// Create a mock file for testing.
	f, err := os.CreateTemp("", "test*.csv")
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	_, err = f.WriteString(`unix,symbol,open,high,low,close
1611902340,BTC/USD,38000.0,39000.0,37000.0,38500.0
1611902440,BTC/USD,38300.0,39300.0,37500.0,39000.0
1611902540,BTC/USD,38600.0,39600.0,37800.0,38500.0
`)
	if err != nil {
		t.Fatalf("failed to write to test file: %v", err)
	}
	err = f.Close()
	if err != nil {
		return
	}

	// Set up the test request.
	fileBody := new(bytes.Buffer)
	writer := multipart.NewWriter(fileBody)
	part, err := writer.CreateFormFile("files", filepath.Base(f.Name()))
	if err != nil {
		t.Fatalf("failed to create multipart form: %v", err)
	}
	_, err = io.Copy(part, f)
	if err != nil {
		t.Fatalf("failed to copy file into multipart form: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return
	}

	// Make the test request.
	req, err := http.NewRequest("POST", "/post", fileBody)
	if err != nil {
		t.Fatalf("failed to create test request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	// Check the response.
	if resp.Code != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, resp.Code)
	}
}
