package models

// ElevationDistance zips elevation changes with distances along route
type ElevationDistance struct {
	Elevation float32 `json:"elevation"`
	Distance float32 `json:"distance"`
}