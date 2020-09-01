package handlers

import (
	"testing"
	"net/http/httptest"

	"blinktag.com/bikesy-wrapper/lib"
)

func TestHealthHandlerOk(t *testing.T) {
	handler := HealthHandler{
		logger: lib.TestLogger(t),
	}
	request := httptest.NewRequest("GET", "http://localhost:8888/health", nil)
	writer := httptest.NewRecorder()
	handler.writeOK(writer, request)
	response := writer.Result()
	if response.StatusCode != 200 {
		t.Errorf("HTTP Response code should be 200 but is %v", response.StatusCode)
	}
}

func TestHealthHandlerNotOk(t *testing.T) {
	handler := HealthHandler{
		logger: lib.TestLogger(t),
	}
	writer := httptest.NewRecorder()
	handler.handleError(500, "an error", writer)
}

func TestHandlerNotNull(t *testing.T) {
	handler := NewHealthHandler(lib.TestLogger(t))
	h := handler.Handler()
	if h == nil {
		t.Errorf("Handler should not be nil")
	}
}