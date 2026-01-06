package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mehrnoosh-hk/devnorth-back/internal/delivery/http/dto"
	"github.com/mehrnoosh-hk/devnorth-back/internal/domain"
)

// Writer handles HTTP responses with structured logging
type Writer struct {
	logger *slog.Logger
}

// NewWriter creates a new response writer with the given logger
func NewWriter(logger *slog.Logger) (*Writer, error) {
	if logger == nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidDependencies, "logger can not be nil")
	}
	return &Writer{
		logger: logger,
	}, nil
}

// JSON sends a JSON response with the given status code
// Only accepts types that implement JSONSerializable to ensure compile-time safety
func (rw *Writer) JSON(w http.ResponseWriter, statusCode int, payload dto.JSONSerializable) {
	w.Header().Set("Content-Type", "application/json")

	if payload != nil {
		// 1. Marshal FIRST (in memory, no network communication yet)
		data, err := json.Marshal(payload)
		if err != nil {
			// 2. Encoding failed? Send 500 BEFORE any body
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"internal_server_error","message":"Failed to encode response"}`))
			return
		}

		// 3. Only send status AFTER we know encoding succeeded
		w.WriteHeader(statusCode)
		w.Write(data) // Now write the valid data
	}
}

// Error sends an error response with appropriate HTTP status code
func (rw *Writer) Error(w http.ResponseWriter, err error) {
	var statusCode int
	var errorCode string
	var message string

	// Map domain errors to HTTP status codes
	switch {
	case errors.Is(err, domain.ErrEmailAlreadyExists):
		statusCode = http.StatusConflict
		errorCode = "email_already_exists"
		message = "An account with this email already exists"

	case errors.Is(err, domain.ErrInvalidCredentials):
		statusCode = http.StatusUnauthorized
		errorCode = "invalid_credentials"
		message = "Invalid email or password"

	case errors.Is(err, domain.ErrInvalidEmail):
		statusCode = http.StatusBadRequest
		errorCode = "invalid_email"
		message = "Invalid email format"

	case errors.Is(err, domain.ErrInvalidPassword):
		statusCode = http.StatusBadRequest
		errorCode = "invalid_password"
		message = "Password must be at least 8 characters"

	case errors.Is(err, domain.ErrInvalidToken):
		statusCode = http.StatusUnauthorized
		errorCode = "invalid_token"
		message = "Invalid or expired token"

	default:
		// Check if it's a validation error
		var validationErr dto.ValidationError
		if errors.As(err, &validationErr) {
			statusCode = http.StatusBadRequest
			errorCode = "validation_error"
			message = validationErr.Error()
		} else {
			// Internal server error for unknown errors
			rw.logger.Error("Unhandled error in response", "error", err)
			statusCode = http.StatusInternalServerError
			errorCode = "internal_server_error"
			message = "An unexpected error occurred"
		}
	}

	resp := dto.ErrorResponse{
		Error:   errorCode,
		Message: message,
	}

	rw.JSON(w, statusCode, resp)
}

// Success sends a successful response
func (rw *Writer) Success(w http.ResponseWriter, payload dto.JSONSerializable) {
	rw.JSON(w, http.StatusOK, payload)
}

// Created sends a created response (201)
func (rw *Writer) Created(w http.ResponseWriter, payload dto.JSONSerializable) {
	rw.JSON(w, http.StatusCreated, payload)
}
