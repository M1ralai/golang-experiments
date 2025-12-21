package http

import (
	"encoding/json"
	"net/http"
)

type HealthHandler struct{}

func NewHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("server is started and running well")
}
