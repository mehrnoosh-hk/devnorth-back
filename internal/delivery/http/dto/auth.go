package dto

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents a successful authentication response
type AuthResponse struct {
	Token   string  `json:"token,omitempty"`
	User    UserDTO `json:"user"`
	Message string  `json:"message,omitempty"`
}

// validateAuthRequest is a shared helper for validating email and password
func validateAuthRequest(email, password string) error {
	if email == "" {
		return ErrFieldRequired("email")
	}
	if password == "" {
		return ErrFieldRequired("password")
	}
	return nil
}

// Validate performs basic validation on RegisterRequest
func (r *RegisterRequest) Validate() error {
	return validateAuthRequest(r.Email, r.Password)
}

// Validate performs basic validation on LoginRequest
func (r *LoginRequest) Validate() error {
	return validateAuthRequest(r.Email, r.Password)
}
