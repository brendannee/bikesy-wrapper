package services

import (
	"github.com/gomodule/redigo/redis"
	"blinktag.com/bikesy-wrapper/config"
	"log"
	"strconv"
	"math"
)

type ElevationService interface {
	GetElevations(nodes []int) ([]float64, error)
}

// ElevationServiceImpl implements ElevationService
// Gets elevation from redis
type ElevationServiceImpl struct {
	redisURL string
}

// NewElevationService sets redis connection and returns ElevationService
func NewElevationService(config *config.Configuration, logger *log.Logger) (ElevationService, error) {
	return &ElevationServiceImpl{
		redisURL: config.Redis.URL,
	}, nil
}


func (s *ElevationServiceImpl) GetElevations(nodes []int) ([]float64, error) {
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
	var elevations []float64
	for _, eString := range v {
		if eString == "" {
			// a value of -1 represents an unknown elevation
			elevations = append(elevations, -1)
		} else {
			eFloat, parseErr := strconv.ParseFloat(eString, 32)
			if parseErr != nil {
				return nil, parseErr
			}
			elevations = append(elevations, math.Round(eFloat*100)/100)
		}
	}
	return elevations, nil
}