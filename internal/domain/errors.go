package domain

import "errors"

// Domain-level errors for business rule violations
var (
	// ErrEmailAlreadyExists is returned when attempting to register with an email that already exists
	ErrEmailAlreadyExists = errors.New("email already exists")

	// ErrInvalidEmail is returned when the email format is invalid
	ErrInvalidEmail = errors.New("invalid email format")

	// ErrInvalidPassword is returned when the password doesn't meet requirements
	ErrInvalidPassword = errors.New("invalid password")

	// ErrUserNotFound is returned when a user cannot be found
	ErrUserNotFound = errors.New("user not found")

	// ErrInvalidCredentials is returned when login credentials are incorrect
	// Generic error to avoid revealing whether email exists (security best practice)
	ErrInvalidCredentials = errors.New("invalid email or password")

	// ErrInvalidToken is returned when a JWT token is invalid or expired
	ErrInvalidToken = errors.New("invalid or expired token")
)
