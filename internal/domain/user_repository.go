package domain

import "context"

// UserRepository defines the contract for user data access
// This interface belongs to the domain layer, defining what operations are needed
// The implementation will be in the repository layer
type UserRepository interface {
	// Create creates a new user in the system
	Create(ctx context.Context, email, hashedPassword string, role UserRole) (*User, error)

	// GetByEmail retrieves a user by their email address
	// Returns domain.ErrUserNotFound if the user doesn't exist
	GetByEmail(ctx context.Context, email string) (*User, error)
}
