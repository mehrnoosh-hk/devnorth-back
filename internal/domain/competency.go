package domain

import "time"

// Competency represents a competency in the domain layer
// This is the core business entity, independent of database implementation
type Competency struct {
	ID          int32
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// HasDescription checks if the competency has a description
func (c *Competency) HasDescription() bool {
	return c.Description != ""
}
