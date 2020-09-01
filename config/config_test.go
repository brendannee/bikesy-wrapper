package config

import (
	"log"
	"os"

	"testing"
)

// testLogger can't be imported from test lib because of circular import
func testLogger(t *testing.T) *log.Logger {
	return log.New(testWriter{t}, "test", log.LstdFlags)
}

type testWriter struct {
	t *testing.T
}

func (tw testWriter) Write(p []byte) (n int, err error) {
	tw.t.Log(string(p))
	return len(p), nil
}

func TestConfigShouldLoadDefaultPort(t *testing.T) {
	logger := testLogger(t)
	os.Setenv("CONFIG", "development.yaml")
	config, err := LoadConfig(logger)
	if err != nil {
		t.Errorf("Development config should load without error %v", err)
	}
	if config.Application.Port != "8888" {
		t.Errorf("If no port specified should be 8888")
	}
}

func TestConfigShouldLoadCustomPort(t *testing.T) {
	logger := testLogger(t)
	os.Setenv("CONFIG", "development.yaml")
	os.Setenv("PORT", "1234")
	config, err := LoadConfig(logger)
	if err != nil {
		t.Errorf("Development config should load without error %v", err)
	}
	if config.Application.Port != "1234" {
		t.Errorf("Config should allow specification of custom port")
	}
}

func TestConfigBadPath(t *testing.T) {
	logger := testLogger(t)
	os.Setenv("CONFIG", "")
	_, err := LoadConfig(logger)
	if err == nil {
		t.Error("Bad config path should error")
	}
}

func TestConfigBadYAML(t *testing.T) {
	logger := testLogger(t)
	os.Setenv("CONFIG", "test_bad_config.yaml")
	_, err := LoadConfig(logger)
	if err == nil {
		t.Error("Bad yaml should error")
	}
}
