package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
	"go.uber.org/config"
)

// Application ...
type Application struct {
	Name string
	Port string
}

// Profile provides routing for each safety preference ("standard", etc.) to correct OSRM instance
type Profile struct {
	Host string
}

// Profiles ...
type Profiles struct {
	HLowSLow Profile `yaml:"h_low_s_low"`
	HLowSMed Profile `yaml:"h_low_s_med"`
	HLowSHigh Profile `yaml:"h_low_s_high"`
	HMedSLow Profile `yaml:"h_med_s_low"`
	HMedSMed Profile `yaml:"h_med_s_med"`
	HMedSHigh Profile `yaml:"h_med_s_high"`
	HHighSLow Profile `yaml:"h_high_s_low"`
	HHighSMed Profile `yaml:"h_high_s_med"`
	HHighSHigh Profile `yaml:"h_high_s_high"`
}

// Osrm includes all profiles
type Osrm struct {
	Profiles Profiles
}

// Redis stores connection info for elevation data
type Redis struct {
	URL string
}

// Configuration ...
type Configuration struct {
	Application Application
	Osrm Osrm
	Redis Redis
}

// LoadConfig reads config.yaml for non secret things, env variables (.env) for others
func LoadConfig(logger *log.Logger) (*Configuration, error) {
	var c Configuration
	// ignore any errors - .env won't exist in production and is instead env var
	godotenv.Load()

	path := os.Getenv("CONFIG")
	port := os.Getenv("PORT")
	cfg, err := config.NewYAML(config.File(path))
	if err != nil {
		return &c, err
	}

	if err := cfg.Get("").Populate(&c); err != nil {
		return &c, err
	}

	// use as secret only
	c.Redis.URL = os.Getenv("REDIS_URL")
	c.Application.Port = port

	return &c, nil
}
