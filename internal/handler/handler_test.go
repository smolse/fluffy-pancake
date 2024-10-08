package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/smolse/fluffy-pancake/internal/datastores"
	"github.com/smolse/fluffy-pancake/internal/service"
)

func TestGetRisk_BadId(t *testing.T) {
	// Create a new handler with a mock service
	svc := &service.RiskService{}
	h := NewHandler(svc)

	// Create a new response recorder
	w := httptest.NewRecorder()

	// Call the handler function with an invalid risk ID
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "bad-uuid-0000-0000-0000-000000000000"})
	h.GetRisk(c)

	// Check that the response status code is 400
	if w.Code != http.StatusBadRequest {
		t.Errorf("unexpected status code: got %d, want %d", w.Code, http.StatusBadRequest)
	}

	// Check that the response body contains the expected error message
	var respBody map[string]string
	json.NewDecoder(w.Body).Decode(&respBody)
	if respBody["error"] != "invalid UUID format" {
		t.Errorf("unexpected error message: got %s, want %s", respBody["error"], "invalid UUID format")
	}
}

func TestGetRisk_NotFound(t *testing.T) {
	// Create a new handler with a mock service
	svc := service.NewRiskService(&datastores.SyncMapDataStore{})
	h := NewHandler(svc)

	// Create a new response recorder
	w := httptest.NewRecorder()

	// Call the handler function with a non-existing risk ID
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "00000000-0000-0000-0000-000000000000"})
	h.GetRisk(c)

	// Check that the response status code is 404
	if w.Code != http.StatusNotFound {
		t.Errorf("unexpected status code: got %d, want %d", w.Code, http.StatusNotFound)
	}

	// Check that the response body contains the expected error message
	var respBody map[string]string
	json.NewDecoder(w.Body).Decode(&respBody)
	if respBody["error"] != "risk 00000000-0000-0000-0000-000000000000 was not found" {
		t.Errorf("unexpected error message: got %s, want %s", respBody["error"], "risk 00000000-0000-0000-0000-000000000000 was not found")
	}
}

func TestCreateRisk_InvalidState(t *testing.T) {
	// Create a new handler with a mock service
	svc := &service.RiskService{}
	h := NewHandler(svc)

	// Create a new request with an invalid risk state
	reqBody := []byte(`{"state": "invalid"}`)
	req, _ := http.NewRequest("POST", "/v1/risks", bytes.NewBuffer(reqBody))

	// Create a new response recorder
	w := httptest.NewRecorder()

	// Call the handler function
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	h.CreateRisk(c)

	// Check that the response status code is 400
	if w.Code != http.StatusBadRequest {
		t.Errorf("unexpected status code: got %d, want %d", w.Code, http.StatusBadRequest)
	}

	// Check that the response body contains the expected error message
	var respBody map[string]string
	json.NewDecoder(w.Body).Decode(&respBody)
	if respBody["error"] != "invalid risk state: invalid" {
		t.Errorf("unexpected error message: got %s, want %s", respBody["error"], "invalid risk state: invalid")
	}
}
