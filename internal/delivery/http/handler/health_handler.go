package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mehrnoosh-hk/devnorth-back/internal/delivery/http/dto"
	"github.com/mehrnoosh-hk/devnorth-back/internal/delivery/http/response"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	responseWriter *response.Writer
}

// NewHealthHandler creates a new health handler instance
func NewHealthHandler(responseWriter *response.Writer) (*HealthHandler, error) {
	if responseWriter == nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidDependencies, "responseWriter can not be nil")
	}
	return &HealthHandler{
		responseWriter: responseWriter,
	}, nil
}

// Check handles health check requests
// GET /api/v1/health
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	resp := dto.HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
	}

	h.responseWriter.Success(w, resp)
}
