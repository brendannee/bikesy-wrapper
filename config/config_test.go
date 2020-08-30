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

func TestConfigShouldLoad(t *testing.T) {
	logger := testLogger(t)
	os.Setenv("MATCHINGCONFIG", "test.yaml")
	_, err := LoadConfig(logger)
	if err != nil {
		t.Errorf("Development config should load without error %v", err)
	}
}

func TestConfigBadPath(t *testing.T) {
	logger := testLogger(t)
	os.Setenv("MATCHINGCONFIG", "")
	_, err := LoadConfig(logger)
	if err == nil {
		t.Error("Bad config path should error")
	}
}

func TestConfigBadYAML(t *testing.T) {
	logger := testLogger(t)
	os.Setenv("MATCHINGCONFIG", "test_bad_config.yaml")
	_, err := LoadConfig(logger)
	if err == nil {
		t.Error("Bad yaml should error")
	}
}
