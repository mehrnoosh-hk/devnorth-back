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

// ToCompetencyDTO converts domain.Competency to CompetencyDTO
func ToCompetencyDTO(competency *domain.Competency, l *slog.Logger) (dto.CompetencyDTO, error) {
	if competency == nil {
		l.Error("Attempt to convert nil domain competency to DTO")
		return dto.CompetencyDTO{}, fmt.Errorf("cannot convert nil domain competency to DTO")
	}

	return dto.CompetencyDTO{
		ID:          competency.ID,
		Name:        competency.Name,
		Description: competency.Description,
		CreatedAt:   competency.CreatedAt,
		UpdatedAt:   competency.UpdatedAt,
	}, nil
}

// ToCompetencyDTOs converts a slice of domain.Competency to a slice of CompetencyDTO
func ToCompetencyDTOs(competencies []*domain.Competency, l *slog.Logger) ([]dto.CompetencyDTO, error) {
	if competencies == nil {
		return []dto.CompetencyDTO{}, nil
	}

	dtos := make([]dto.CompetencyDTO, len(competencies))
	for i, competency := range competencies {
		dto, err := ToCompetencyDTO(competency, l)
		if err != nil {
			return nil, err
		}
		dtos[i] = dto
	}

	return dtos, nil
}
