package lib

import (
	"log"
	"testing"

	"blinktag.com/bikesy-wrapper/config"
)

// TestLogger ...
func TestLogger(t *testing.T) *log.Logger {
	return log.New(testWriter{t}, "test", log.LstdFlags)
}

type testWriter struct {
	t *testing.T
}

func (tw testWriter) Write(p []byte) (n int, err error) {
	tw.t.Log(string(p))
	return len(p), nil
}

// TestConfig ...
func TestConfig(name string) *config.Configuration {
	return &config.Configuration{
		Application: config.Application{
			Name: name,
		},
	}
} 
