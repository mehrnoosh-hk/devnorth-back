package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mehrnoosh-hk/devnorth-back/internal/delivery/http/dto"
	"github.com/mehrnoosh-hk/devnorth-back/internal/delivery/http/response"
	"github.com/mehrnoosh-hk/devnorth-back/internal/domain"
)

// CompetencyHandler handles competency-related HTTP requests
type CompetencyHandler struct {
	competencyUseCase domain.CompetencyUseCase
	logger            *slog.Logger
	responseWriter    *response.Writer
}

// NewCompetencyHandler creates a new competency handler instance
func NewCompetencyHandler(competencyUseCase domain.CompetencyUseCase, logger *slog.Logger, responseWriter *response.Writer) (*CompetencyHandler, error) {
	// Check if dependencies are nil
	if competencyUseCase == nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidDependencies, "competencyUseCase can not be nil")
	}
	if logger == nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidDependencies, "logger can not be nil")
	}
	if responseWriter == nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidDependencies, "responseWriter can not be nil")
	}
	return &CompetencyHandler{
		competencyUseCase: competencyUseCase,
		logger:            logger,
		responseWriter:    responseWriter,
	}, nil
}

// Create handles competency creation requests
// POST /api/v1/competencies
// HTTP Status Codes:
//   - 201 Created: Competency created successfully
//   - 400 Bad Request: Validation errors (invalid name)
//   - 409 Conflict: Competency name already exists
//   - 500 Internal Server Error: Unexpected errors
func (h *CompetencyHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCompetencyRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Failed to decode create competency request", "error", err)
		h.responseWriter.Error(w, ErrInvalidJSON)
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		h.logger.Warn("Create competency request validation failed", "error", err)
		h.responseWriter.Error(w, err)
		return
	}

	// Call use case
	competency, err := h.competencyUseCase.Create(r.Context(), req.Name, req.Description)
	if err != nil {
		h.logger.Error("Failed to create competency", "error", err)
		h.responseWriter.Error(w, err)
		return
	}

	// Build response
	competencyDTO, err := ToCompetencyDTO(competency, h.logger)
	if err != nil {
		h.responseWriter.Error(w, err)
		return
	}

	h.logger.Info("Competency created successfully", "competency_id", competency.ID, "name", competency.Name)
	h.responseWriter.Created(w, competencyDTO)
}

// GetByID handles get competency by ID requests
// GET /api/v1/competencies/{id}
// HTTP Status Codes:
//   - 200 OK: Competency retrieved successfully
//   - 400 Bad Request: Invalid ID format
//   - 404 Not Found: Competency not found
//   - 500 Internal Server Error: Unexpected errors
func (h *CompetencyHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL parameter
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("Invalid competency ID format", "id", idStr, "error", err)
		h.responseWriter.Error(w, dto.ValidationError{
			Field:   "id",
			Message: "invalid ID format",
		})
		return
	}

	// Call use case
	competency, err := h.competencyUseCase.GetByID(r.Context(), int32(id))
	if err != nil {
		h.logger.Error("Failed to get competency by ID", "id", id, "error", err)
		h.responseWriter.Error(w, err)
		return
	}

	// Build response
	competencyDTO, err := ToCompetencyDTO(competency, h.logger)
	if err != nil {
		h.responseWriter.Error(w, err)
		return
	}

	h.logger.Info("Competency retrieved successfully", "competency_id", competency.ID)
	h.responseWriter.Success(w, competencyDTO)
}

// GetAll handles get all competencies requests
// GET /api/v1/competencies
// HTTP Status Codes:
//   - 200 OK: Competencies retrieved successfully
//   - 500 Internal Server Error: Unexpected errors
func (h *CompetencyHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Call use case
	competencies, err := h.competencyUseCase.GetAll(r.Context())
	if err != nil {
		h.logger.Error("Failed to get all competencies", "error", err)
		h.responseWriter.Error(w, err)
		return
	}

	// Build response
	competencyDTOs, err := ToCompetencyDTOs(competencies, h.logger)
	if err != nil {
		h.responseWriter.Error(w, err)
		return
	}

	resp := dto.CompetenciesResponse{
		Competencies: competencyDTOs,
		Count:        len(competencyDTOs),
	}

	h.logger.Info("Competencies retrieved successfully", "count", len(competencies))
	h.responseWriter.Success(w, resp)
}

// UpdateDescription handles update competency description requests
// PATCH /api/v1/competencies/{id}/description
// HTTP Status Codes:
//   - 200 OK: Description updated successfully
//   - 400 Bad Request: Invalid ID format
//   - 404 Not Found: Competency not found
//   - 500 Internal Server Error: Unexpected errors
func (h *CompetencyHandler) UpdateDescription(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL parameter
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("Invalid competency ID format", "id", idStr, "error", err)
		h.responseWriter.Error(w, dto.ValidationError{
			Field:   "id",
			Message: "invalid ID format",
		})
		return
	}

	// Parse request body
	var req dto.UpdateCompetencyDescriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Failed to decode update competency description request", "error", err)
		h.responseWriter.Error(w, ErrInvalidJSON)
		return
	}

	// Validate request (currently no validation needed, but keeping for consistency)
	if err := req.Validate(); err != nil {
		h.logger.Warn("Update competency description request validation failed", "error", err)
		h.responseWriter.Error(w, err)
		return
	}

	// Call use case
	competency, err := h.competencyUseCase.UpdateDescription(r.Context(), int32(id), req.Description)
	if err != nil {
		h.logger.Error("Failed to update competency description", "id", id, "error", err)
		h.responseWriter.Error(w, err)
		return
	}

	// Build response
	competencyDTO, err := ToCompetencyDTO(competency, h.logger)
	if err != nil {
		h.responseWriter.Error(w, err)
		return
	}

	h.logger.Info("Competency description updated successfully", "competency_id", competency.ID)
	h.responseWriter.Success(w, competencyDTO)
}
