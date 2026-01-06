package usecase

import "errors"

var (
	// Dependency errors
	ErrUserRepositoryNil = errors.New("user repository cannot be nil")
	ErrPasswordHasherNil = errors.New("password hasher cannot be nil")
	ErrTokenGeneratorNil = errors.New("token generator cannot be nil")
	ErrLoggerNil         = errors.New("logger cannot be nil")

	// Operation errors
	ErrCheckExistingUser = errors.New("failed to check existing user")
	ErrHashPassword      = errors.New("failed to hash password")
	ErrCreateUser        = errors.New("failed to create user")
	ErrGenerateToken     = errors.New("failed to generate token")
	ErrGetUser           = errors.New("failed to get user")
)
