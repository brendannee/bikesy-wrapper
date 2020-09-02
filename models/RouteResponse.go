package models

// Annotation ...
type Annotation struct {
	Distance []float32 `json:"distance"`
	Nodes []int `json:"nodes"`
}

// Intersection ...
type Intersection struct {
	Out int `json:"out"`
	In int `json:"int"`
	Entry []bool `json:"entry"`
	Bearings []int `json:"bearings"`
	Location []float32 `json:"location"`
}

// Maneuver ...
type Maneuver struct {
	BearingAfter int `json:"bearing_after"`
	Type string `json:"type"`
	Modifier string `json:"modifier"`
	BearingBefore int `json:"bearing_before"`
	Location []float32 `json:"location"`
}

// Step ...
type Step struct {
	Intersections []Intersection `json:"intersections"`
	DrivingSide string `json:"driving_side"`
	Geometry string `json:"geometry"`
	Mode string `json:"mode"`
	Duration float32 `json:"duration"`
	Maneuver Maneuver `json:"maneuver"`
	Weight float32 `json:"weight"`
	Distance float32 `json:"distance"`
	Name string `json:"name"`
	Pronunciation string `json:"pronunciation"`
}

// Leg ...
type Leg struct {
	Annotation Annotation `json:"annotation"`
	Steps []Step `json:"steps"`
}

// Route ...
type Route struct {
	Geometry string `json:"geometry"`
	Legs []Leg `json:"legs"`
}

// RouteResponse is fields of interest from bikesy api
type RouteResponse struct {
	Routes []Route `json:"routes"`
}
