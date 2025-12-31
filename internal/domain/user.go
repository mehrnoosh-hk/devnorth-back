package domain

import "time"

// UserRole represents the user's role in the system
type UserRole string

const (
	UserRoleUSER  UserRole = "USER"
	UserRoleADMIN UserRole = "ADMIN"
)

// User represents a user in the domain layer
// This is the core business entity, independent of database implementation
type User struct {
	ID             int32
	Email          string
	HashedPassword string
	Role           UserRole
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// IsAdmin checks if the user has admin privileges
func (u *User) IsAdmin() bool {
	return u.Role == UserRoleADMIN
}

// IsUser checks if the user has regular user privileges
func (u *User) IsUser() bool {
	return u.Role == UserRoleUSER
}
