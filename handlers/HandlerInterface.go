package handlers

import (
	"net/http"
)

// Handler interface for all service handlers, which must export
// http.Handler to be handled by Mux
type Handler interface {
	Handler() (http.Handler)
	handleError(statusCode int, errorMsg string, w http.ResponseWriter)
	handleOK(response interface{}, w http.ResponseWriter)
}