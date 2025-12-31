package security

import "errors"

var (
	// Configuration errors
	ErrKeyIsTooShort          = errors.New("key is too short")
	ErrDurationMustBePositive = errors.New("duration must be positive")
	ErrLoggerCanNotBeNil      = errors.New("logger can not be nil")
	ErrInvalidBcryptCost      = errors.New("invalid bcrypt cost")

	// Token errors
	ErrFailedToSignToken       = errors.New("failed to sign token")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrUnexpectedKeyID         = errors.New("unexpected key ID")
	ErrUserCanNotBeNil         = errors.New("user can not be nil")

	// Password errors
	ErrHashPassword    = errors.New("failed to hash password")
	ErrComparePassword = errors.New("failed to compare passwords")
)
