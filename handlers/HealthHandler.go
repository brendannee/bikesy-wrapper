package handlers

import (
	"log"
	"encoding/json"
	"net/http"

	"blinktag.com/bikesy-wrapper/models"
)

// HealthHandler returns Google Polyline from Packen Service stored in Firebase
type HealthHandler struct {
	logger *log.Logger
}

// NewHealthHandler ...
func NewHealthHandler(logger *log.Logger) (Handler) {
	return &HealthHandler {
		logger: logger,
	}
}

func (h* HealthHandler) handleError(statusCode int, errorMsg string, w http.ResponseWriter) {
	http.Error(w, errorMsg, statusCode)
}

func (h* HealthHandler) writeOK(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := models.HealthHandlerResponse{Status: "ok"}
	h.handleOK(response, w)
}

func (h* HealthHandler) handleOK(response interface{}, w http.ResponseWriter) {
	json.NewEncoder(w).Encode(response)
}

// Handler implements Handler interface
func (h *HealthHandler) Handler()  (http.Handler) {
	h.logger.Print("Executing health check handler.")
	return http.HandlerFunc(h.writeOK)
}