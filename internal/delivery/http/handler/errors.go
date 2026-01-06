package handler

import (
	"errors"

	"github.com/mehrnoosh-hk/devnorth-back/internal/delivery/http/dto"
)

var (

	// ErrInvalidDependencies is returned when the dependencies are nil
	ErrInvalidDependencies = errors.New("invalid dependencies")

	// ErrInvalidJSON is a response returned when the request body cannot be decoded
	ErrInvalidJSON = dto.ValidationError{
		Field:   "body",
		Message: "invalid JSON format",
	}
)
