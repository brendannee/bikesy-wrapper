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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// won't error here
	args, _ := url.ParseQuery(r.URL.RawQuery)

	// TO DO: define other profiles
    var profileType services.ProfileType

    var lat1 string
    var lng1 string
    var lat2 string
    var lng2 string
    var hills string
    var safety string

    var val []string
    var ok bool

    if val, ok = args["hills"]; ok {
        hills = val[0]
    } else {
        h.handleError(400, "Request requires hill tolerance", w)
        return 
    }
    if val, ok = args["safety"]; ok {
        safety = val[0]
    } else {
        h.handleError(400, "Request requires safety tolerance", w)
        return 
    }

    // check for valid combinations

    if (hills == "low" && safety == "low") {
        profileType = services.ProfileTypeHLowSLow
    } else if (hills == "low" && safety == "med") {
        profileType = services.ProfileTypeHLowSMed
    } else if (hills == "low" && safety == "high") {
        profileType = services.ProfileTypeHLowSHigh
    } else if (hills == "med" && safety == "low") {
        profileType = services.ProfileTypeHMedSLow
    } else if (hills == "med" && safety == "med") {
        profileType = services.ProfileTypeHMedSMed
    } else if (hills == "med" && safety == "high") {
        profileType = services.ProfileTypeHMedSHigh
    } else if (hills == "high" && safety == "low") {
        profileType = services.ProfileTypeHHighSLow
    } else if (hills == "high" && safety == "med") {
        profileType = services.ProfileTypeHHighSMed
    } else if (hills == "high" && safety == "high") {
        profileType = services.ProfileTypeHHighSHigh
    } else {
        h.handleError(400, "Bad combination of hill tolerance and safety", w)
        return 
    }

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
    elevationProfile, err := h.elevationService.GetElevationsAndDistances(nodes, resp.Routes[0].Legs[0].Annotation.Distance)
    if (err != nil) {
        // treat any redis issues as 500
        h.logger.Printf("Error parsing elevation data %v", err)
        h.handleError(500, err.Error(), w)
        return
    }

    bikesyResponse := models.BikesyResponse{
        Geometry: resp.Routes[0].Geometry,
        ElevationProfile: elevationProfile,
        Steps: resp.Routes[0].Legs[0].Steps,
    }
	h.handleOK(bikesyResponse, w)
}

// Handler implements Handler interface
func (h *BikesyHandler) Handler()  (http.Handler) {
	h.logger.Print("Executing route handler.")
	return http.HandlerFunc(h.handleRouteRequest)
}