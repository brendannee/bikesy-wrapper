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

func validRouteResponse() models.RouteResponse {
	var mockRoutes []models.Route
	var mockLegs []models.Leg
	var mockDistance []float32
	var mockNodes []int
	mockDistance = append(mockDistance, 2)
	mockNodes = append(mockNodes, 1)
	mockAnnotation := models.Annotation{
		Distance: mockDistance,
		Nodes: mockNodes,
	}
	mockLeg := models.Leg{
		Annotation: mockAnnotation,
	}
	mockLegs = append(mockLegs, mockLeg)
	mockRoute := models.Route{
		Geometry: "geometry",
		Legs: mockLegs,
	}
	mockRoutes = append(mockRoutes, mockRoute)
	return models.RouteResponse{
		Routes: mockRoutes,
	}
}

func TestBikesyHandlerRequiredParams(t *testing.T) {
	handler := NewBikesyHandler(lib.TestLogger(t), &mocks.RouteService{}, &mocks.ElevationService{})
	response := handler.Handler()
	if response == nil {
		t.Errorf("HTTP handler func should not be nil")
	}
}

func TestBikesyHandlerOk(t *testing.T) {
	mockRouteService := &mocks.RouteService{}
	mockElevationService := &mocks.ElevationService{}
	routeResponse := validRouteResponse()
	var mockElevation []models.ElevationDistance

	mockRouteService.On("GetBikeRoute", "1", "2", "3", "4").Return(routeResponse, nil)
	mockRouteService.On("SetProfile", services.ProfileType("HILLS_LOW_SAFETY_LOW")).Return()
	mockElevationService.On("GetElevationsAndDistances", routeResponse.Routes[0].Legs[0].Annotation.Nodes, routeResponse.Routes[0].Legs[0].Annotation.Distance).Return(mockElevation, nil)
	handler := BikesyHandler{
		logger: lib.TestLogger(t),
		routeService: mockRouteService,
		elevationService: mockElevationService,
	}
	request := httptest.NewRequest("GET", "http://localhost:8888/route?hills=low&safety=low&lat1=1&lng1=2&lat2=3&lng2=4", nil)
	writer := httptest.NewRecorder()
	handler.handleRouteRequest(writer, request)
	response := writer.Result()
	if response.StatusCode != 200 {
		t.Errorf("HTTP Response code should be 200 but is %v", response.StatusCode)
	}
}

func TestBikesyHandlerBadElevationResponse(t *testing.T) {
	mockRouteService := &mocks.RouteService{}
	mockElevationService := &mocks.ElevationService{}
	routeResponse := validRouteResponse()

	mockRouteService.On("GetBikeRoute", "1", "2", "3", "4").Return(routeResponse, nil)
	mockRouteService.On("SetProfile", services.ProfileType("HILLS_LOW_SAFETY_LOW")).Return()
	mockElevationService.On("GetElevationsAndDistances", routeResponse.Routes[0].Legs[0].Annotation.Nodes, routeResponse.Routes[0].Legs[0].Annotation.Distance).Return(nil, errors.New("something happened"))
	handler := BikesyHandler{
		logger: lib.TestLogger(t),
		routeService: mockRouteService,
		elevationService: mockElevationService,
	}
	request := httptest.NewRequest("GET", "http://localhost:8888/route?hills=low&safety=low&lat1=1&lng1=2&lat2=3&lng2=4", nil)
	writer := httptest.NewRecorder()
	handler.handleRouteRequest(writer, request)
	response := writer.Result()
	if response.StatusCode != 500 {
		t.Errorf("HTTP Response code should be 500 but is %v", response.StatusCode)
	}
}

func TestBikesyHandlerBadOSRMRequest(t *testing.T) {
	mockRouteService := &mocks.RouteService{}
	mockRouteService.On("GetBikeRoute", "1", "2", "3", "4").Return(models.RouteResponse{}, errors.New("something happened"))
	mockRouteService.On("SetProfile", services.ProfileType("HILLS_LOW_SAFETY_LOW")).Return()
	handler := BikesyHandler{
		logger: lib.TestLogger(t),
		routeService: mockRouteService,
	}
	request := httptest.NewRequest("GET", "http://localhost:8888/route?hills=low&safety=low&lat1=1&lng1=2&lat2=3&lng2=4", nil)
	writer := httptest.NewRecorder()
	handler.handleRouteRequest(writer, request)
	response := writer.Result()
	if response.StatusCode != 500 {
		t.Errorf("HTTP Response code should be 500 but is %v", response.StatusCode)
	}
}

func TestBikesyHandlerBadOSRMResponseTooManyRoutes(t *testing.T) {
	invalidResponse := validRouteResponse()
	invalidResponse.Routes = append(invalidResponse.Routes, models.Route{})
	mockRouteService := &mocks.RouteService{}
	mockRouteService.On("GetBikeRoute", "1", "2", "3", "4").Return(invalidResponse, nil)
	mockRouteService.On("SetProfile", services.ProfileType("HILLS_LOW_SAFETY_LOW")).Return()
	handler := BikesyHandler{
		logger: lib.TestLogger(t),
		routeService: mockRouteService,
	}
	request := httptest.NewRequest("GET", "http://localhost:8888/route?hills=low&safety=low&lat1=1&lng1=2&lat2=3&lng2=4", nil)
	writer := httptest.NewRecorder()
	handler.handleRouteRequest(writer, request)
	response := writer.Result()
	if response.StatusCode != 500 {
		t.Errorf("HTTP Response code should be 500 but is %v", response.StatusCode)
	}
}

func TestBikesyHandlerBadOSRMResponseTooManyLegs(t *testing.T) {
	invalidResponse := validRouteResponse()
	invalidResponse.Routes[0].Legs = append(invalidResponse.Routes[0].Legs, models.Leg{})
	mockRouteService := &mocks.RouteService{}
	mockRouteService.On("GetBikeRoute", "1", "2", "3", "4").Return(invalidResponse, nil)
	mockRouteService.On("SetProfile", services.ProfileType("HILLS_LOW_SAFETY_LOW")).Return()
	handler := BikesyHandler{
		logger: lib.TestLogger(t),
		routeService: mockRouteService,
	}
	request := httptest.NewRequest("GET", "http://localhost:8888/route?hills=low&safety=low&lat1=1&lng1=2&lat2=3&lng2=4", nil)
	writer := httptest.NewRecorder()
	handler.handleRouteRequest(writer, request)
	response := writer.Result()
	if response.StatusCode != 500 {
		t.Errorf("HTTP Response code should be 500 but is %v", response.StatusCode)
	}
}

func TestBikesyHandlerAllArgsRequired(t *testing.T) {
	mockRouteService := &mocks.RouteService{}
	handler := BikesyHandler{
		logger: lib.TestLogger(t),
		routeService: mockRouteService,
	}
	// Missing Lng2
	request := httptest.NewRequest("GET", "http://localhost:8888/route?hills=low&safety=low&lat1=1&lng1=2&lat2=3", nil)
	writer := httptest.NewRecorder()
	handler.handleRouteRequest(writer, request)
	response := writer.Result()
	if response.StatusCode != 400 {
		t.Errorf("HTTP Response code should be 400 when missing arguments but is %v", response.StatusCode)
	}

	// Missing Lat2
	request = httptest.NewRequest("GET", "http://localhost:8888/route?hills=low&safety=low&lat1=1&lng1=2&lng2=3", nil)
	writer = httptest.NewRecorder()
	handler.handleRouteRequest(writer, request)
	response = writer.Result()
	if response.StatusCode != 400 {
		t.Errorf("HTTP Response code should be 400 when missing arguments but is %v", response.StatusCode)
	}

	// Missing Lng1
	request = httptest.NewRequest("GET", "http://localhost:8888/route?hills=low&safety=low&lat1=1&lat2=2&lng2=3", nil)
	writer = httptest.NewRecorder()
	handler.handleRouteRequest(writer, request)
	response = writer.Result()
	if response.StatusCode != 400 {
		t.Errorf("HTTP Response code should be 400 when missing arguments but is %v", response.StatusCode)
	}

	// Missing Lat1
	request = httptest.NewRequest("GET", "http://localhost:8888/route?hills=low&safety=low&lng1=1&lat2=2&lng2=3", nil)
	writer = httptest.NewRecorder()
	handler.handleRouteRequest(writer, request)
	response = writer.Result()
	if response.StatusCode != 400 {
		t.Errorf("HTTP Response code should be 400 when missing arguments but is %v", response.StatusCode)
	}
}