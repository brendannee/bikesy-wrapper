package models

// BikesyResponse is fields of interest from bikesy api and elevation
type BikesyResponse struct {
	Geometry string `json:"geometry"`
	ElevationProfile []ElevationDistance `json:"elevation_profile"`
	Steps []Step `json:"steps"`
}
