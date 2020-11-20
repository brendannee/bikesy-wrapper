package services

import (
	"github.com/gomodule/redigo/redis"
	"blinktag.com/bikesy-wrapper/config"
	"blinktag.com/bikesy-wrapper/models"
	"strconv"
	"math"
	"fmt"
)

// ElevationService gets elevations for front-end display
type ElevationService interface {
	GetElevationsAndDistances(nodes []int, distances []float32) ([]models.ElevationDistance, error)
}

// ElevationServiceImpl implements ElevationService
// Gets elevation from redis
type ElevationServiceImpl struct {
	redisURL string
}

// NewElevationService sets redis connection and returns ElevationService
func NewElevationService(config *config.Configuration) ElevationService {
	return &ElevationServiceImpl{
		redisURL: config.Redis.URL,
	}
}

// GetElevationsAndDistances ...
func (s *ElevationServiceImpl) GetElevationsAndDistances(nodes []int, distances []float32) ([]models.ElevationDistance, error) {
	fmt.Printf("%v", s.redisURL)
	c, err := redis.DialURL(s.redisURL)
	if err != nil {
	    return nil, err
	}
	defer c.Close()
	c.Send("MULTI")
	for _, n := range nodes {
		c.Send("GET", n)
	}
	v, redisErr := redis.Strings(c.Do("EXEC"))
	if redisErr != nil {
	    return nil, redisErr
	}
	var elevations []models.ElevationDistance
	distance := float32(0)
	for i, eString := range v {
		var elevation float32
		if eString == "" {
			// a value of -1 represents an unknown elevation
			elevation = float32(-1)
		} else {
			eFloat, parseErr := strconv.ParseFloat(eString, 32)
			if parseErr != nil {
				return nil, parseErr
			}
			elevation = float32(math.Round(eFloat * 100) / 100)
		}
		elevations = append(elevations, models.ElevationDistance{
			Elevation: elevation,
			Distance: distance,
		})
		if i < len(distances) {
			distance = distance + distances[i]
		}
	}

	return elevations, nil
}