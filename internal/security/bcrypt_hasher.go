package security

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/mehrnoosh-hk/devnorth-back/internal/domain"
)

// bcryptHasher implements domain.PasswordHasher using bcrypt algorithm
type bcryptHasher struct {
	cost int
}

// NewBcryptHasher creates a pointer to a new bcryptHasher and error which satisfies domain.PasswordHasher interface
// If error happens, NewBcryptHasher returns nil and error
// The caller must check the nil pointer
// cost determines the computational cost of hashing (higher = more secure but slower)
// Default cost is bcrypt.DefaultCost (10)
func NewBcryptHasher(cost int) (*bcryptHasher, error) {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return nil, fmt.Errorf("invalid bcrypt cost %d, must be between %d and %d", 
			cost, bcrypt.MinCost, bcrypt.MaxCost)
	}
	return &bcryptHasher{cost: cost}, nil
}

// Hash generates a bcrypt hash from a plain text password and error which satisfies domain.PasswordHasher interface
// If error happens, Hash returns "" and error
func (h *bcryptHasher) Hash(password string) (string, error) {
	// Make sure no DoS attack
	if len(password) > 72 {
		return "", domain.ErrInvalidCredentials
	}
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

// Compare verifies if a plain text password matches a bcrypt hashed password
func (h *bcryptHasher) Compare(hashedPassword, plainPassword string) error {
	// Make sure no DoS attack
	if len(plainPassword) > 72 {
		return domain.ErrInvalidCredentials
	}
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return domain.ErrInvalidCredentials
		}
		return fmt.Errorf("failed to compare passwords: %w", err)
	}
	return nil
}
