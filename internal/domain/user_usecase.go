package domain

import "context"

// UserUseCase defines the contract for user-related business operations
// This interface belongs to the domain layer and will be implemented by the use case layer
type UserUseCase interface {
	// Register creates a new user account with the provided email and password
	// Returns the created user (without password) or an error if registration fails
	// Possible errors: ErrEmailAlreadyExists, ErrInvalidEmail, ErrInvalidPassword
	Register(ctx context.Context, email, password string) (*User, error)

	// Login authenticates a user with email and password
	// Returns a JWT token and the user (without password) if authentication succeeds
	// Possible errors: ErrInvalidCredentials
	Login(ctx context.Context, email, password string) (string, *User, error)
}
