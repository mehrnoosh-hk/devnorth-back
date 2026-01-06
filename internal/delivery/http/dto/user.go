package dto

import "time"

// UserDTO represents user data in API responses
// Excludes sensitive information like hashed password
type UserDTO struct {
	ID        int32     `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// Implement JSONSerializable
func (UserDTO) isJSONSerializable() {}
