package dto

import (
	"log"
	"time"

	"github.com/mehrnoosh-hk/devnorth-back/internal/domain"
)

// AuthResponse represents a successful authentication response
type AuthResponse struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}

// UserDTO represents user data in API responses
// Excludes sensitive information like hashed password
type UserDTO struct {
	ID        int32     `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// ToUserDTO converts domain.User to UserDTO
func ToUserDTO(user *domain.User) UserDTO {
	if user == nil {
		log.Println("Warning: Attempted to convert nil user to DTO")
		return UserDTO{}
	}

	return UserDTO{
		ID:        user.ID,
		Email:     user.Email,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
	}
}
