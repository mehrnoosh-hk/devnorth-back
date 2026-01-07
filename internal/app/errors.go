package app

import "errors"

var (
	// Initialization errors
	ErrInitDB                = errors.New("failed to initialize database connection")
	ErrInitPasswordHasher    = errors.New("failed to initialize password hasher")
	ErrInitTokenGenerator    = errors.New("failed to initialize token generator")
	ErrInitUserUseCase       = errors.New("failed to initialize user use case")
	ErrInitCompetencyUseCase = errors.New("failed to initialize competency use case")
)
