package handler

import (
	"github.com/mehrnoosh-hk/devnorth-back/internal/delivery/http/dto"
)

var (
	// ErrInvalidJSON is returned when the request body cannot be decoded
	ErrInvalidJSON = dto.ValidationError{
		Field:   "body",
		Message: "invalid JSON format",
	}
)
