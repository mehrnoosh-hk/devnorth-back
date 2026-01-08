package dto

import "time"

// CreateCompetencyRequest represents the request to create a new competency
type CreateCompetencyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// UpdateCompetencyDescriptionRequest represents the request to update a competency's description
type UpdateCompetencyDescriptionRequest struct {
	Description string `json:"description"`
}

// CompetencyDTO represents competency data in API responses
type CompetencyDTO struct {
	ID          int32     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CompetenciesResponse represents a list of competencies in API responses
type CompetenciesResponse struct {
	Competencies []CompetencyDTO `json:"competencies"`
	Count        int             `json:"count"`
}

// Validate performs basic validation on CreateCompetencyRequest
func (r *CreateCompetencyRequest) Validate() error {
	if r.Name == "" {
		return ErrFieldRequired("name")
	}
	return nil
}

// Validate performs basic validation on UpdateCompetencyDescriptionRequest
// Description can be empty (to clear it), so no validation needed
func (r *UpdateCompetencyDescriptionRequest) Validate() error {
	return nil
}

// Implement JSONSerializable for all competency DTOs
func (CreateCompetencyRequest) isJSONSerializable()              {}
func (UpdateCompetencyDescriptionRequest) isJSONSerializable()   {}
func (CompetencyDTO) isJSONSerializable()                        {}
func (CompetenciesResponse) isJSONSerializable()                 {}
