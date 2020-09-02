package models

// BikesyResponse is fields of interest from bikesy api and elevation
type BikesyResponse struct {
	Geometry string
	Elevation []float64
	Distance []float32
	Steps []Step
}
