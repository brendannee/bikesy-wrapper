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

// ProfileTypeHLowSLow is low hills low safety
const ProfileTypeHLowSLow = ProfileType("HILLS_LOW_SAFETY_LOW")

// ProfileTypeHLowSMed ...
const ProfileTypeHLowSMed = ProfileType("HILLS_LOW_SAFETY_MED")

// ProfileTypeHLowSHigh ...
const ProfileTypeHLowSHigh = ProfileType("HILLS_LOW_SAFETY_HIGH")

// ProfileTypeHMedSLow ...
const ProfileTypeHMedSLow = ProfileType("HILLS_MED_SAFETY_LOW")

// ProfileTypeHMedSMed ...
const ProfileTypeHMedSMed = ProfileType("HILLS_MED_SAFETY_MED")

// ProfileTypeHMedSHigh ...
const ProfileTypeHMedSHigh = ProfileType("HILLS_MED_SAFETY_HIGH")

// ProfileTypeHHighSLow ...
const ProfileTypeHHighSLow = ProfileType("HILLS_HIGH_SAFETY_LOW")

// ProfileTypeHHighSMed ...
const ProfileTypeHHighSMed = ProfileType("HILLS_HIGH_SAFETY_MED")

// ProfileTypeHHighSHigh ...
const ProfileTypeHHighSHigh = ProfileType("HILLS_HIGH_SAFETY_HIGH")


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
	var urlBase string
	if (s.profile == ProfileTypeHLowSLow) {
		urlBase = s.config.Osrm.Profiles.HLowSLow.Host
	} else if (s.profile == ProfileTypeHLowSMed) {
		urlBase = s.config.Osrm.Profiles.HLowSMed.Host
	} else if (s.profile == ProfileTypeHLowSHigh) {
		urlBase = s.config.Osrm.Profiles.HLowSHigh.Host
	} else if (s.profile == ProfileTypeHMedSLow) {
		urlBase = s.config.Osrm.Profiles.HMedSLow.Host
	} else if (s.profile == ProfileTypeHMedSMed) {
		urlBase = s.config.Osrm.Profiles.HMedSMed.Host
	} else if (s.profile == ProfileTypeHMedSHigh) {
		urlBase = s.config.Osrm.Profiles.HMedSHigh.Host
	} else if (s.profile == ProfileTypeHHighSLow) {
		urlBase = s.config.Osrm.Profiles.HHighSLow.Host
	} else if (s.profile == ProfileTypeHHighSMed) {
		urlBase = s.config.Osrm.Profiles.HHighSMed.Host
	} else if (s.profile == ProfileTypeHHighSHigh) {
		urlBase = s.config.Osrm.Profiles.HHighSHigh.Host
	} else {
		return response, errors.New("profile is not supported")
	}

	s.logger.Printf("Requesting profile: %s, URL: %s", s.profile, urlBase)

	// Get response from matching server
	resp, err := http.Get(fmt.Sprintf("%s%s,%s;%s,%s?steps=true&annotations=true&overview=full", urlBase, lng1, lat1, lng2, lat2))
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
