package services

import (
	"testing"
	"errors"
	"net/http"
	"net/http/httptest"
	"fmt"

	"blinktag.com/bikesy-wrapper/lib"
	"blinktag.com/bikesy-wrapper/config"
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
    return 0, errors.New("test error")
}

func TestRouteServiceStandardOnly(t *testing.T) {
	service := NewRouteService(lib.TestConfig("test"), lib.TestLogger(t))
	_, err := service.GetBikeRoute("0", "0", "0", "0")
	if err == nil {
		t.Errorf("Route service currently only supports standard route")
	}
}

func TestRouteServiceOk(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		response := "{some json}"
		rw.Write([]byte(response))
	}))
	// Close the server when test finishes
	defer server.Close()
	cfg := lib.TestConfig("test")
	cfg.Osrm = config.Osrm {
		Profiles: config.Profiles{
			HLowSLow: config.Profile{
				Host: fmt.Sprintf("%s/", server.URL),
			},
		},
	}
	s := NewRouteService(cfg, lib.TestLogger(t))
	s.SetProfile(ProfileTypeHLowSLow)
	_, err := s.GetBikeRoute("0", "0", "0", "0")
	if err != nil {
		t.Errorf("Route service should not error but did with %s", err.Error())
	}
}

func TestRouteServiceNotOkBadOSRMResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		response := "{some json}"
		rw.WriteHeader(500)
		rw.Write([]byte(response))
	}))
	// Close the server when test finishes
	defer server.Close()
	cfg := lib.TestConfig("test")
	cfg.Osrm = config.Osrm {
		Profiles: config.Profiles{
			HLowSLow: config.Profile{
				Host: fmt.Sprintf("%s/", server.URL),
			},
		},
	}
	s := NewRouteService(cfg, lib.TestLogger(t))
	s.SetProfile(ProfileTypeHLowSLow)
	_, err := s.GetBikeRoute("0", "0", "0", "0")
	if err == nil {
		t.Errorf("Route service should error if OSRM status code != 200")
	}
}

func TestRouteServiceBadOsrmRequest(t *testing.T) {
	cfg := lib.TestConfig("test")
	cfg.Osrm = config.Osrm {
		Profiles: config.Profiles{
			HLowSLow: config.Profile{
				Host: "A Bad URL",
			},
		},
	}
	// server.URL = fmt.Sprintf("%s/%s,%s;%s,%s?steps=false&annotations=true", server.URL, "0", "0", "0", "0")
	s := NewRouteService(cfg, lib.TestLogger(t))
	s.SetProfile(ProfileTypeHLowSLow)
	_, err := s.GetBikeRoute("0", "0", "0", "0")
	if err == nil {
		t.Errorf("Route service should error with bad OSRM URL")
	}
}
