package handlers

import (
	"errors"
	"testing"
	"net/http/httptest"

	"blinktag.com/bikesy-wrapper/models"
	"blinktag.com/bikesy-wrapper/services/mocks"
	"blinktag.com/bikesy-wrapper/services"
	"blinktag.com/bikesy-wrapper/lib"
)

func TestBikesyHandlerRequiredParams(t *testing.T) {
	handler := NewBikesyHandler(lib.TestLogger(t), &mocks.RouteService{}, &mocks.ElevationService{})
	response := handler.Handler()
	if response == nil {
		t.Errorf("HTTP handler func should not be nil")
	}
}

func TestBikesyHandlerOk(t *testing.T) {
	mockService := &mocks.RouteService{}
	mockService.On("GetBikeRoute", "1", "2", "3", "4").Return(models.RouteResponse{}, nil)
	mockService.On("SetProfile", services.ProfileType("STANDARD")).Return()
	handler := BikesyHandler{
		logger: lib.TestLogger(t),
		routeService: mockService,
	}
	request := httptest.NewRequest("GET", "http://localhost:8888/route?lat1=1&lng1=2&lat2=3&lng2=4", nil)
	writer := httptest.NewRecorder()
	handler.handleRouteRequest(writer, request)
	response := writer.Result()
	if response.StatusCode != 200 {
		t.Errorf("HTTP Response code should be 200 but is %v", response.StatusCode)
	}
}

func TestBikesyHandlerBadOSRMRequest(t *testing.T) {
	mockService := &mocks.RouteService{}
	mockService.On("GetBikeRoute", "1", "2", "3", "4").Return(models.RouteResponse{}, errors.New("something happened"))
	mockService.On("SetProfile", services.ProfileType("STANDARD")).Return()
	handler := BikesyHandler{
		logger: lib.TestLogger(t),
		routeService: mockService,
	}
	request := httptest.NewRequest("GET", "http://localhost:8888/route?lat1=1&lng1=2&lat2=3&lng2=4", nil)
	writer := httptest.NewRecorder()
	handler.handleRouteRequest(writer, request)
	response := writer.Result()
	if response.StatusCode != 500 {
		t.Errorf("HTTP Response code should be 500 but is %v", response.StatusCode)
	}
}

func TestBikesyHandlerAllArgsRequired(t *testing.T) {
	mockService := &mocks.RouteService{}
	handler := BikesyHandler{
		logger: lib.TestLogger(t),
		routeService: mockService,
	}
	// Missing Lng2
	request := httptest.NewRequest("GET", "http://localhost:8888/route?lat1=1&lng1=2&lat2=3", nil)
	writer := httptest.NewRecorder()
	handler.handleRouteRequest(writer, request)
	response := writer.Result()
	if response.StatusCode != 400 {
		t.Errorf("HTTP Response code should be 400 when missing arguments but is %v", response.StatusCode)
	}

	// Missing Lat2
	request = httptest.NewRequest("GET", "http://localhost:8888/route?lat1=1&lng1=2&lng2=3", nil)
	writer = httptest.NewRecorder()
	handler.handleRouteRequest(writer, request)
	response = writer.Result()
	if response.StatusCode != 400 {
		t.Errorf("HTTP Response code should be 400 when missing arguments but is %v", response.StatusCode)
	}

	// Missing Lng1
	request = httptest.NewRequest("GET", "http://localhost:8888/route?lat1=1&lat2=2&lng2=3", nil)
	writer = httptest.NewRecorder()
	handler.handleRouteRequest(writer, request)
	response = writer.Result()
	if response.StatusCode != 400 {
		t.Errorf("HTTP Response code should be 400 when missing arguments but is %v", response.StatusCode)
	}

	// Missing Lat1
	request = httptest.NewRequest("GET", "http://localhost:8888/route?lng1=1&lat2=2&lng2=3", nil)
	writer = httptest.NewRecorder()
	handler.handleRouteRequest(writer, request)
	response = writer.Result()
	if response.StatusCode != 400 {
		t.Errorf("HTTP Response code should be 400 when missing arguments but is %v", response.StatusCode)
	}
}