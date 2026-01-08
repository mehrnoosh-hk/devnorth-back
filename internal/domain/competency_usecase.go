package domain

import "context"

// CompetencyUseCase defines the contract for competency-related business operations
// This interface belongs to the domain layer and will be implemented by the use case layer
type CompetencyUseCase interface {
	// Create creates a new competency with the provided name and description
	// Returns the created competency or an error if creation fails
	// Possible errors: ErrCompetencyAlreadyExists, ErrInvalidCompetencyName
	Create(ctx context.Context, name, description string) (*Competency, error)

	// GetByID retrieves a competency by its ID
	// Returns domain.ErrCompetencyNotFound if the competency doesn't exist
	GetByID(ctx context.Context, id int32) (*Competency, error)

	// GetByName retrieves a competency by its name
	// Returns domain.ErrCompetencyNotFound if the competency doesn't exist
	GetByName(ctx context.Context, name string) (*Competency, error)

	// GetAll retrieves all competencies
	GetAll(ctx context.Context) ([]*Competency, error)

	// UpdateDescription updates the description of a competency
	// Returns domain.ErrCompetencyNotFound if the competency doesn't exist
	UpdateDescription(ctx context.Context, id int32, description string) (*Competency, error)
}
