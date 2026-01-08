package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/mehrnoosh-hk/devnorth-back/internal/domain"
)

// competencyUseCase implements domain.CompetencyUseCase
// It orchestrates competency-related business operations using repository
type competencyUseCase struct {
	competencyRepo domain.CompetencyRepository
	logger         *slog.Logger
}

// NewCompetencyUseCase creates a new competency use case instance
// Dependencies are injected following the Dependency Inversion Principle
func NewCompetencyUseCase(
	competencyRepo domain.CompetencyRepository,
	logger *slog.Logger,
) (domain.CompetencyUseCase, error) {
	// Nil-check the injected dependencies
	if competencyRepo == nil {
		return nil, ErrCompetencyRepositoryNil
	}
	if logger == nil {
		return nil, ErrLoggerNil
	}
	return &competencyUseCase{
		competencyRepo: competencyRepo,
		logger:         logger,
	}, nil
}

// Create creates a new competency
// Business logic flow:
// 1. Validate name (basic validation for POC)
// 2. Normalize name (trim spaces)
// 3. Check if competency with same name already exists
// 4. Create competency in repository
func (uc *competencyUseCase) Create(ctx context.Context, name, description string) (*domain.Competency, error) {
	// Normalize name by trimming spaces
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)

	// Step 1: Validate name (basic validation for POC)
	if err := uc.validateName(name); err != nil {
		uc.logger.Error("failed to validate competency name", "error", err)
		return nil, err
	}

	// Step 2: Check if competency already exists
	existingCompetency, err := uc.competencyRepo.GetByName(ctx, name)
	if err != nil && !errors.Is(err, domain.ErrCompetencyNotFound) {
		uc.logger.Error("failed to check existing competency", "error", err)
		return nil, fmt.Errorf("%w: %w", ErrCheckExistingCompetency, err)
	}
	if existingCompetency != nil {
		uc.logger.Error("competency already exists", "name", name)
		return nil, domain.ErrCompetencyAlreadyExists
	}

	// Step 3: Create competency
	competency, err := uc.competencyRepo.Create(ctx, name, description)
	if err != nil {
		uc.logger.Error("failed to create competency", "error", err)
		return nil, fmt.Errorf("%w: %w", ErrCreateCompetency, err)
	}

	uc.logger.Info("competency created successfully", "competency_id", competency.ID, "name", competency.Name)
	return competency, nil
}

// GetByID retrieves a competency by its ID
func (uc *competencyUseCase) GetByID(ctx context.Context, id int32) (*domain.Competency, error) {
	competency, err := uc.competencyRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrCompetencyNotFound) {
			uc.logger.Info("competency not found", "id", id)
			return nil, domain.ErrCompetencyNotFound
		}
		uc.logger.Error("failed to get competency by ID", "error", err, "id", id)
		return nil, fmt.Errorf("%w: %w", ErrGetCompetency, err)
	}

	return competency, nil
}

// GetByName retrieves a competency by its name
func (uc *competencyUseCase) GetByName(ctx context.Context, name string) (*domain.Competency, error) {
	// Normalize name
	name = strings.TrimSpace(name)

	competency, err := uc.competencyRepo.GetByName(ctx, name)
	if err != nil {
		if errors.Is(err, domain.ErrCompetencyNotFound) {
			uc.logger.Info("competency not found", "name", name)
			return nil, domain.ErrCompetencyNotFound
		}
		uc.logger.Error("failed to get competency by name", "error", err, "name", name)
		return nil, fmt.Errorf("%w: %w", ErrGetCompetency, err)
	}

	return competency, nil
}

// GetAll retrieves all competencies
func (uc *competencyUseCase) GetAll(ctx context.Context) ([]*domain.Competency, error) {
	competencies, err := uc.competencyRepo.GetAll(ctx)
	if err != nil {
		uc.logger.Error("failed to get all competencies", "error", err)
		return nil, fmt.Errorf("%w: %w", ErrGetCompetencies, err)
	}

	uc.logger.Info("competencies retrieved successfully", "count", len(competencies))
	return competencies, nil
}

// UpdateDescription updates the description of a competency
// Business logic flow:
// 1. Validate that competency exists (implicitly done by repository)
// 2. Update description in repository
func (uc *competencyUseCase) UpdateDescription(ctx context.Context, id int32, description string) (*domain.Competency, error) {
	// Normalize description
	description = strings.TrimSpace(description)

	competency, err := uc.competencyRepo.UpdateDescription(ctx, id, description)
	if err != nil {
		if errors.Is(err, domain.ErrCompetencyNotFound) {
			uc.logger.Info("competency not found for update", "id", id)
			return nil, domain.ErrCompetencyNotFound
		}
		uc.logger.Error("failed to update competency description", "error", err, "id", id)
		return nil, fmt.Errorf("%w: %w", ErrUpdateCompetency, err)
	}

	uc.logger.Info("competency description updated successfully", "competency_id", competency.ID)
	return competency, nil
}

// validateName performs basic name validation
// For POC: minimal validation - just check it's not empty and has reasonable length
// Production: add more sophisticated validation rules
func (uc *competencyUseCase) validateName(name string) error {
	name = strings.TrimSpace(name)

	// Basic validation for POC
	if name == "" {
		return domain.ErrInvalidCompetencyName
	}

	// Check minimum length
	if len(name) < 2 {
		return domain.ErrInvalidCompetencyName
	}

	// Check maximum length (reasonable limit for a competency name)
	if len(name) > 100 {
		return domain.ErrInvalidCompetencyName
	}

	// TODO for production:
	// - Check for invalid characters
	// - Check against reserved names
	// - More sophisticated validation rules

	return nil
}
