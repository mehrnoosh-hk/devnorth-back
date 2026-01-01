package domain

// PasswordHasher defines the contract for password hashing operations
// This interface belongs to the domain layer, allowing the use case to depend on abstraction
// rather than concrete implementations (Dependency Inversion Principle)
type PasswordHasher interface {
	// Hash takes a plain text password and returns a hashed version
	Hash(password string) (string, error)

	// Compare verifies if a plain text password matches a hashed password
	Compare(hashedPassword, plainPassword string) error
}
