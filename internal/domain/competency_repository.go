package domain

import "context"

// CompetencyRepository defines the contract for competency data access
// This interface belongs to the domain layer, defining what operations are needed
// The implementation will be in the repository layer
type CompetencyRepository interface {
	// Create creates a new competency in the system
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
