package usecase

import "errors"

var (
	// Dependency errors
	ErrUserRepositoryNil       = errors.New("user repository cannot be nil")
	ErrCompetencyRepositoryNil = errors.New("competency repository cannot be nil")
	ErrPasswordHasherNil       = errors.New("password hasher cannot be nil")
	ErrTokenGeneratorNil       = errors.New("token generator cannot be nil")
	ErrLoggerNil               = errors.New("logger cannot be nil")

	// User operation errors
	ErrCheckExistingUser = errors.New("failed to check existing user")
	ErrHashPassword      = errors.New("failed to hash password")
	ErrCreateUser        = errors.New("failed to create user")
	ErrGenerateToken     = errors.New("failed to generate token")
	ErrGetUser           = errors.New("failed to get user")

	// Competency operation errors
	ErrCheckExistingCompetency = errors.New("failed to check existing competency")
	ErrCreateCompetency        = errors.New("failed to create competency")
	ErrGetCompetency           = errors.New("failed to get competency")
	ErrGetCompetencies         = errors.New("failed to get competencies")
	ErrUpdateCompetency        = errors.New("failed to update competency")
)
