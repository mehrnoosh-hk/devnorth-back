package repository

import "errors"

var (
	// Dependency errors
	ErrPoolNil   = errors.New("database pool cannot be nil")
	ErrLoggerNil = errors.New("logger cannot be nil")

	// Competency repository errors
	ErrCreateCompetencyFailed            = errors.New("failed to create competency")
	ErrGetCompetencyByIDFailed           = errors.New("failed to get competency by ID")
	ErrGetCompetencyByNameFailed         = errors.New("failed to get competency by name")
	ErrGetAllCompetenciesFailed          = errors.New("failed to get all competencies")
	ErrUpdateCompetencyDescriptionFailed = errors.New("failed to update competency description")
)
