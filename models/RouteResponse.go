package models

// Annotation ...
type Annotation struct {
	Distance []float32 `json:"distance"`
	Nodes []int `json:"nodes"`
}

// Leg ...
type Leg struct {
	Annotation Annotation `json:"annotation"`
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
