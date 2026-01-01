package handler

import (
	"net/http"
	"time"

	"github.com/mehrnoosh-hk/devnorth-back/internal/delivery/http/dto"
	"github.com/mehrnoosh-hk/devnorth-back/internal/delivery/http/response"
)

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new health handler instance
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Check handles health check requests
// GET /api/v1/health
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	resp := dto.HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
	}

	response.Success(w, resp)
}
