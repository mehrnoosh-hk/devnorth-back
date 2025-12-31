package handler

import (
	"fmt"
	"log/slog"

	"github.com/mehrnoosh-hk/devnorth-back/internal/delivery/http/dto"
	"github.com/mehrnoosh-hk/devnorth-back/internal/domain"
)

// ToUserDTO converts domain.User to UserDTO
func ToUserDTO(user *domain.User, l *slog.Logger) (dto.UserDTO, error) {
	if user == nil {
		l.Error("Attempt to convert nil domain user to DTO")
		return dto.UserDTO{}, fmt.Errorf("cannot convert nil domain user to DTO")
	}

	return dto.UserDTO{
		ID:        user.ID,
		Email:     user.Email,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
	}, nil
}
