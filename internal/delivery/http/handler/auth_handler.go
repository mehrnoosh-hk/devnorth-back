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
	userUseCase    domain.UserUseCase
	logger         *slog.Logger
	responseWriter *response.Writer
}

// NewAuthHandler creates a new auth handler instance
func NewAuthHandler(userUseCase domain.UserUseCase, logger *slog.Logger, responseWriter *response.Writer) *AuthHandler {
	return &AuthHandler{
		userUseCase:    userUseCase,
		logger:         logger,
		responseWriter: responseWriter,
	}
}

// Register handles user registration requests with automatic login
// POST /api/v1/auth/register
// HTTP Status Codes:
//   - 201 Created: User registered successfully (with token if auto-login succeeded)
//   - 201 Created: User registered but auto-login failed (without token, with message)
//   - 400 Bad Request: Validation errors (invalid email/password format)
//   - 409 Conflict: Email already exists
//   - 500 Internal Server Error: Unexpected errors
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Failed to decode register request", "error", err)
		h.responseWriter.Error(w, ErrInvalidJSON)
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		h.logger.Warn("Register request validation failed", "error", err)
		h.responseWriter.Error(w, err)
		return
	}

	// Step 1: Register the user
	user, err := h.userUseCase.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Error("Failed to register user", "error", err)
		h.responseWriter.Error(w, err)
		return
	}

	h.logger.Info("User registered successfully", "user_id", user.ID)

	// Step 2: Attempt automatic login
	token, _, err := h.userUseCase.Login(r.Context(), req.Email, req.Password)
	userDTO, dtoErr := ToUserDTO(user, h.logger)
	if dtoErr != nil {
		h.responseWriter.Error(w, dtoErr)
		return
	}

	if err != nil {
		// Registration succeeded but auto-login failed
		// Still return 201 Created (resource was created) but without token
		h.logger.Warn("User registered but auto-login failed", "error", err, "user_id", user.ID)
		h.responseWriter.Created(w, dto.AuthResponse{
			User:    userDTO,
			Message: "Account created successfully. Please try logging in.",
		})
		return
	}

	// Both registration and login succeeded
	h.logger.Info("User registered and logged in successfully", "user_id", user.ID)
	h.responseWriter.Created(w, dto.AuthResponse{
		Token: token,
		User:  userDTO,
	})
}

// Login handles user login requests
// POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Failed to decode login request", "error", err)
		h.responseWriter.Error(w, ErrInvalidJSON)
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		h.logger.Warn("Login request validation failed", "error", err)
		h.responseWriter.Error(w, err)
		return
	}

	// Call use case
	token, user, err := h.userUseCase.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Warn("Login failed",
			"error", err,
		)
		h.responseWriter.Error(w, err)
		return
	}

	// Build response
	dtoUser, err := ToUserDTO(user, h.logger)
	if err != nil {
		h.responseWriter.Error(w, err)
		return
	}
	resp := dto.AuthResponse{
		Token: token,
		User:  dtoUser,
	}

	h.logger.Info("User logged in successfully",
		"user_id", user.ID,
	)

	h.responseWriter.Success(w, resp)
}
