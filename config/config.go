package config

import (
	"log"
	"os"

	"go.uber.org/config"
)

// Application ...
type Application struct {
	Name string
	Port string
}

// DB ...

// Configuration ...
type Configuration struct {
	Application Application
}

// LoadConfig reads development.yaml for now
func LoadConfig(logger *log.Logger) (*Configuration, error) {
	var c Configuration
	path := os.Getenv("CONFIG")
	port := os.Getenv("PORT")
	cfg, err := config.NewYAML(config.File(path))
	if err != nil {
		return &c, err
	}

	if err := cfg.Get("").Populate(&c); err != nil {
		return &c, err
	}

	if port != "" {
		logger.Printf("Using custom port %v", port)
		c.Application.Port = port
	}
	return &c, nil
}
