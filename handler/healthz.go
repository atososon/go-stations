package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
)

// A HealthzHandler implements health check endpoint.
type HealthzHandler struct{}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewHealthzHandler() *HealthzHandler {
	return &HealthzHandler{}
}

// ServeHTTP implements http.Handler interface.
func (h *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// set response header
	w.Header().Set("Content-Type", "application/json")

	// create response
	resp := model.HealthzResponse{
		Message: "OK",
	}

	// encode response
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println("healthz: failed to encode response, err =", err)
	}
}