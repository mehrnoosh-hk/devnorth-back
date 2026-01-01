package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/mehrnoosh-hk/devnorth-back/internal/delivery/http/dto"
	"github.com/mehrnoosh-hk/devnorth-back/internal/delivery/http/response"
	"github.com/mehrnoosh-hk/devnorth-back/internal/domain"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	userUseCase domain.UserUseCase
	logger      *slog.Logger
}

// NewAuthHandler creates a new auth handler instance
func NewAuthHandler(userUseCase domain.UserUseCase, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		userUseCase: userUseCase,
		logger:      logger,
	}
}

// Register handles user registration requests
// POST /api/v1/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Failed to decode register request", "error", err)
		response.Error(w, dto.ValidationError{
			Field:   "body",
			Message: "invalid JSON format",
		})
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		h.logger.Warn("Register request validation failed", "error", err)
		response.Error(w, err)
		return
	}

	// Call use case
	user, err := h.userUseCase.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Error("Failed to register user",
			"email", req.Email,
			"error", err,
		)
		response.Error(w, err)
		return
	}

	// Generate token for automatic login after registration
	token, _, err := h.userUseCase.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		// User registered but token generation failed
		// This shouldn't happen, but handle gracefully
		h.logger.Error("Failed to generate token after registration",
			"user_id", user.ID,
			"error", err,
		)
		response.Error(w, err)
		return
	}

	// Build response
	resp := dto.AuthResponse{
		Token: token,
		User:  dto.ToUserDTO(user),
	}

	h.logger.Info("User registered successfully",
		"user_id", user.ID,
		"email", user.Email,
	)

	response.Created(w, resp)
}

// Login handles user login requests
// POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Failed to decode login request", "error", err)
		response.Error(w, dto.ValidationError{
			Field:   "body",
			Message: "invalid JSON format",
		})
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		h.logger.Warn("Login request validation failed", "error", err)
		response.Error(w, err)
		return
	}

	// Call use case
	token, user, err := h.userUseCase.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Warn("Login failed",
			"email", req.Email,
			"error", err,
		)
		response.Error(w, err)
		return
	}

	// Build response
	resp := dto.AuthResponse{
		Token: token,
		User:  dto.ToUserDTO(user),
	}

	h.logger.Info("User logged in successfully",
		"user_id", user.ID,
		"email", user.Email,
	)

	response.Success(w, resp)
}
