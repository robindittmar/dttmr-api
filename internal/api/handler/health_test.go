package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/robindittmar/dttmr-api/internal/api/handler"
)

func TestHealthHandler_Success(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	handler.HealthHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("HealthHandler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	expected := `{"status":"ok"}`
	if rr.Body.String() != expected {
		t.Errorf("HealthHandler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
