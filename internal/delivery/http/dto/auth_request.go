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

// validateAuthRequest is a shared helper for validating email and password
// For POC: minimal validation, production should use validation library
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
