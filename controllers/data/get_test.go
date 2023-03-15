package data_test

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"ohcl/controllers/data"
	"ohcl/controllers/outputForms"
	"testing"
)

func TestGetEndpoint(t *testing.T) {
	// Set up the test server.
	r := gin.Default()
	r.GET("/data", data.Get)

	// Seed the database with some test data.
	//database.DB.Create(&ohcl.Model{
	//	Symbol:    "BTC/USD",
	//	Unix:      1611902340,
	//	Open:      38000.0,
	//	High:      39000.0,
	//	Low:       37000.0,
	//	Close:     38500.0,
	//	CreatedAt: time.Now().Unix(),
	//})

	// Set up the test request.
	req, err := http.NewRequest("GET", "/data?page=1&per_page=5", nil)
	if err != nil {
		t.Fatalf("failed to create test request: %v", err)
	}
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	// Check the response.
	if resp.Code != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, resp.Code)
	}
	var state outputForms.State
	if err = json.NewDecoder(resp.Body).Decode(&state); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if !state.Status {
		t.Errorf("expected status to be true; got %v", state.Status)
	}
}
