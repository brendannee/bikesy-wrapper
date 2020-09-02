package handlers

import (
	"log"
	"encoding/json"
	"net/http"
	"net/url"

	"blinktag.com/bikesy-wrapper/services"
    "blinktag.com/bikesy-wrapper/models"
)

// BikesyHandler returns raw response from standard osrm router for now
type BikesyHandler struct {
	logger *log.Logger
	routeService services.RouteService
    elevationService services.ElevationService
}

// NewBikesyHandler ...
func NewBikesyHandler(logger *log.Logger, routeService services.RouteService, elevationService services.ElevationService) (Handler) {
	return &BikesyHandler {
		logger: logger,
		routeService: routeService,
        elevationService: elevationService,
	}
}

func (h* BikesyHandler) handleError(statusCode int, errorMsg string, w http.ResponseWriter) {
	http.Error(w, errorMsg, statusCode)
}

func (h* BikesyHandler) handleOK(response interface{}, w http.ResponseWriter) {
	json.NewEncoder(w).Encode(response)
}

func (h *BikesyHandler) handleRouteRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// won't error here
	args, _ := url.ParseQuery(r.URL.RawQuery)

	// TO DO: define other profiles
    profileType := services.ProfileTypeStandard

    var lat1 string
    var lng1 string
    var lat2 string
    var lng2 string
 
    var val []string
    var ok bool

    if val, ok = args["lat1"]; ok {
    	lat1 = val[0]
    } else {
    	h.handleError(400, "Request requires lat1", w)
    	return
    }
    if val, ok = args["lng1"]; ok {
    	lng1 = val[0]
    } else {
    	h.handleError(400, "Request requires lng1", w)
    	return
    }
    if val, ok = args["lat2"]; ok {
    	lat2 = val[0]
    } else {
    	h.handleError(400, "Request requires lat2", w)
    	return
    }
    if val, ok = args["lng2"]; ok {
    	lng2 = val[0]
    } else {
    	h.handleError(400, "Request requires lng2", w)
    	return
    }
    h.logger.Printf("Received request for %v %v %v %v", lat1, lng1, lat2, lng2)

    h.routeService.SetProfile(profileType)
    // Get response from route service
	resp, err := h.routeService.GetBikeRoute(lat1, lng1, lat2, lng2)
	if (err != nil) {
		// treat any connection issues as 500 for alerts
		h.logger.Printf("Error connecting to osrm service %v", err)
		h.handleError(500, err.Error(), w)
		return
	}
    routes := resp.Routes
    if len(routes) != 1 {
        h.logger.Printf("Osrm data should contain exactly one route")
        h.handleError(500, "Bad response from OSRM server", w)
        return
    }
    legs := routes[0].Legs
    if len(legs) != 1 {
        h.logger.Printf("Osrm data should contain exactly one leg")
        h.handleError(500, "Bad response from OSRM server", w)
        return
    }
    nodes := legs[0].Annotation.Nodes
    elevation, err := h.elevationService.GetElevations(nodes)
    if (err != nil) {
        // treat any redis issues as 500
        h.logger.Printf("Error parsing elevation data %v", err)
        h.handleError(500, err.Error(), w)
        return
    }
    bikesyResponse := models.BikesyResponse{
        Geometry: resp.Routes[0].Geometry,
        Elevation: elevation,
        Distance: resp.Routes[0].Legs[0].Annotation.Distance,
        Steps: resp.Routes[0].Legs[0].Steps,
    }
	h.handleOK(bikesyResponse, w)
}

// Handler implements Handler interface
func (h *BikesyHandler) Handler()  (http.Handler) {
	h.logger.Print("Executing route handler.")
	return http.HandlerFunc(h.handleRouteRequest)
}