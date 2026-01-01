package domain

import "context"

// TokenGenerator defines the contract for JWT token operations
// This interface belongs to the domain layer, allowing the use case to depend on abstraction
// rather than concrete JWT implementations (Dependency Inversion Principle)
type TokenGenerator interface {
	// Generate creates a JWT token for the given user
	// The token contains user claims (ID, email, role), token expiration and kid
	Generate(ctx context.Context, user *User) (string, error)

	// Validate verifies a JWT token and returns the user claims
	// Returns ErrInvalidToken if token is invalid, expired, or malformed
	Validate(ctx context.Context, token string) (*User, error)
}
