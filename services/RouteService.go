package services

import (
	"errors"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"fmt"

	"blinktag.com/bikesy-wrapper/config"
	"blinktag.com/bikesy-wrapper/models"
)

// ProfileType defines how to direct service when specifying desired level of safety
type ProfileType = string

// ProfileTypeStandard is a medium-safe bicycle route
const ProfileTypeStandard = ProfileType("STANDARD")

// RouteService is interface for testing
type RouteService interface {
	GetBikeRoute(lat1 string, lng1 string, lat2 string, lng2 string) (models.RouteResponse, error)
	SetProfile(profile ProfileType)
}

// RouteServiceImpl implements RouteService
// Gets bike route from biksey api
type RouteServiceImpl struct {
	logger *log.Logger
	config *config.Configuration
	profile ProfileType
}

// NewRouteService ...
func NewRouteService(config *config.Configuration, logger *log.Logger) RouteService {
	return &RouteServiceImpl{
		logger: logger,
		config: config,
		profile: ProfileType(""),
	}
}

// SetProfile specifies desired safety level
func (s *RouteServiceImpl) SetProfile(profile ProfileType) {
	s.profile = profile
}

// GetBikeRoute returns OSRM bike route given a safety profile
func (s *RouteServiceImpl) GetBikeRoute(lat1 string, lng1 string, lat2 string, lng2 string) (models.RouteResponse, error) {
	response := models.RouteResponse{}
	if (s.profile != ProfileTypeStandard) {
		return response, errors.New("only supports standard profile for now")
	}
	urlBase := s.config.Osrm.Profiles.Standard.Host
	
	// Get response from matching server
	resp, err := http.Get(fmt.Sprintf("%s%s,%s;%s,%s?steps=true&annotations=true", urlBase, lng1, lat1, lng2, lat2))
	if (err != nil || resp == nil) {
		s.logger.Printf("Error connecting to osrm service %v", err)
		return response, err
	}
	if resp.StatusCode != 200 {
		s.logger.Printf("Received bad response code from osrm service %v", resp.StatusCode)
		return response, errors.New("Bad status code")
	}
	body, bodyErr := ioutil.ReadAll(resp.Body)
	if bodyErr != nil {
		s.logger.Printf("Error reading body from osrm service %v", bodyErr)
		return response, bodyErr
	}
	json.Unmarshal(body, &response)
	return response, nil
}
