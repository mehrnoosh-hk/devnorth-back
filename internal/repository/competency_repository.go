package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mehrnoosh-hk/devnorth-back/db/sqlc"
	"github.com/mehrnoosh-hk/devnorth-back/internal/domain"
)

// competencyRepository implements domain.CompetencyRepository using SQLC
type competencyRepository struct {
	queries *sqlc.Queries
	logger  *slog.Logger
}

// NewCompetencyRepository creates a new instance of CompetencyRepository
func NewCompetencyRepository(pool *pgxpool.Pool, logger *slog.Logger) (domain.CompetencyRepository, error) {
	if pool == nil {
		return nil, ErrPoolNil
	}
	if logger == nil {
		return nil, ErrLoggerNil
	}
	return &competencyRepository{
		queries: sqlc.New(pool),
		logger:  logger,
	}, nil
}

// Create creates a new competency in the database
func (r *competencyRepository) Create(ctx context.Context, name, description string) (*domain.Competency, error) {
	r.logger.Info("creating competency", "name", name)

	params := sqlc.CreateCompetencyParams{
		Name:        name,
		Description: pgtype.Text{String: description, Valid: true},
	}

	sqlcCompetency, err := r.queries.CreateCompetency(ctx, params)
	if err != nil {
		// Check for unique constraint violation (duplicate name)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			r.logger.Warn("duplicate competency name", "name", name)
			return nil, domain.ErrCompetencyAlreadyExists
		}
		r.logger.Error("failed to create competency", "error", err, "name", name)
		return nil, fmt.Errorf("%w: %w", ErrCreateCompetencyFailed, err)
	}

	r.logger.Info("competency created successfully", "competency_id", sqlcCompetency.ID)

	// Convert SQLC model to domain model
	return toDomainCompetency(sqlcCompetency), nil
}

// GetByID retrieves a competency by ID
func (r *competencyRepository) GetByID(ctx context.Context, id int32) (*domain.Competency, error) {
	r.logger.Info("getting competency by ID", "id", id)

	sqlcCompetency, err := r.queries.GetCompetencyByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Info("competency not found", "id", id)
			return nil, domain.ErrCompetencyNotFound
		}
		r.logger.Error("failed to get competency by ID", "error", err, "id", id)
		return nil, fmt.Errorf("%w: %w", ErrGetCompetencyByIDFailed, err)
	}

	r.logger.Info("competency retrieved successfully", "competency_id", sqlcCompetency.ID)

	// Convert SQLC model to domain model
	return toDomainCompetency(sqlcCompetency), nil
}

// GetByName retrieves a competency by name
func (r *competencyRepository) GetByName(ctx context.Context, name string) (*domain.Competency, error) {
	r.logger.Info("getting competency by name", "name", name)

	sqlcCompetency, err := r.queries.GetCompetencyByName(ctx, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Info("competency not found", "name", name)
			return nil, domain.ErrCompetencyNotFound
		}
		r.logger.Error("failed to get competency by name", "error", err, "name", name)
		return nil, fmt.Errorf("%w: %w", ErrGetCompetencyByNameFailed, err)
	}

	r.logger.Info("competency retrieved successfully", "competency_id", sqlcCompetency.ID)

	// Convert SQLC model to domain model
	return toDomainCompetency(sqlcCompetency), nil
}

// GetAll retrieves all competencies ordered by creation date
func (r *competencyRepository) GetAll(ctx context.Context) ([]*domain.Competency, error) {
	r.logger.Info("getting all competencies")

	sqlcCompetencies, err := r.queries.GetAllCompetencies(ctx)
	if err != nil {
		r.logger.Error("failed to get all competencies", "error", err)
		return nil, fmt.Errorf("%w: %w", ErrGetAllCompetenciesFailed, err)
	}

	r.logger.Info("competencies retrieved successfully", "count", len(sqlcCompetencies))

	// Convert SQLC models to domain models
	competencies := make([]*domain.Competency, len(sqlcCompetencies))
	for i, sqlcComp := range sqlcCompetencies {
		competencies[i] = toDomainCompetency(sqlcComp)
	}

	return competencies, nil
}

// UpdateDescription updates the description of a competency
func (r *competencyRepository) UpdateDescription(ctx context.Context, id int32, description string) (*domain.Competency, error) {
	r.logger.Info("updating competency description", "id", id)

	params := sqlc.UpdateCompetencyDescriptionParams{
		ID:          id,
		Description: pgtype.Text{String: description, Valid: true},
	}

	sqlcCompetency, err := r.queries.UpdateCompetencyDescription(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Info("competency not found", "id", id)
			return nil, domain.ErrCompetencyNotFound
		}
		r.logger.Error("failed to update competency description", "error", err, "id", id)
		return nil, fmt.Errorf("%w: %w", ErrUpdateCompetencyDescriptionFailed, err)
	}

	r.logger.Info("competency description updated successfully", "competency_id", sqlcCompetency.ID)

	// Convert SQLC model to domain model
	return toDomainCompetency(sqlcCompetency), nil
}

// toDomainCompetency converts SQLC Competency model to domain Competency model
func toDomainCompetency(sqlcCompetency sqlc.Competency) *domain.Competency {
	var createdAt, updatedAt time.Time
	var description string

	if sqlcCompetency.CreatedAt.Valid {
		createdAt = sqlcCompetency.CreatedAt.Time
	}

	if sqlcCompetency.UpdatedAt.Valid {
		updatedAt = sqlcCompetency.UpdatedAt.Time
	}

	if sqlcCompetency.Description.Valid {
		description = sqlcCompetency.Description.String
	}

	return &domain.Competency{
		ID:          sqlcCompetency.ID,
		Name:        sqlcCompetency.Name,
		Description: description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
