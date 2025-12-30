package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mehrnoosh-hk/devnorth-back/db/sqlc"
	"github.com/mehrnoosh-hk/devnorth-back/internal/domain"
)

// userRepository implements domain.UserRepository using SQLC
type userRepository struct {
	queries *sqlc.Queries
	logger  *slog.Logger
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(pool *pgxpool.Pool, logger *slog.Logger) domain.UserRepository {
	return &userRepository{
		queries: sqlc.New(pool),
		logger:  logger,
	}
}

// Create creates a new user in the database
func (r *userRepository) Create(ctx context.Context, email, hashedPassword string, role domain.UserRole) (*domain.User, error) {
	r.logger.Info("creating user", "role", role)

	params := sqlc.CreateUserParams{
		Email:          email,
		HashedPassword: hashedPassword,
		Role:           sqlc.UserRole(role), // Convert domain.UserRole to sqlc.UserRole
	}

	sqlcUser, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		// Check for unique constraint violation (duplicate email)
		var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				r.logger.Warn("duplicate email", "email", email)
				return nil, domain.ErrEmailAlreadyExists
		}
		r.logger.Error("failed to create user", "error", err, "Role", role)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	r.logger.Info("user created successfully", "user_id", sqlcUser.ID)

	// Convert SQLC model to domain model
	return toDomainUser(sqlcUser), nil
}

// GetByEmail retrieves a user by email address
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.logger.Info("getting user by email")

	sqlcUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.logger.Info("user not found")
			return nil, domain.ErrUserNotFound
		}
		r.logger.Error("failed to get user by email", "error", err)
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	r.logger.Info("user retrieved successfully", "user_id", sqlcUser.ID)

	// Convert SQLC model to domain model
	return toDomainUser(sqlcUser), nil
}

// toDomainUser converts SQLC User model to domain User model
func toDomainUser(sqlcUser sqlc.User) *domain.User {
	var createdAt, updatedAt time.Time

	if sqlcUser.CreatedAt.Valid {
		createdAt = sqlcUser.CreatedAt.Time
	}

	if sqlcUser.UpdatedAt.Valid {
		updatedAt = sqlcUser.UpdatedAt.Time
	}

	return &domain.User{
		ID:             sqlcUser.ID,
		Email:          sqlcUser.Email,
		HashedPassword: sqlcUser.HashedPassword,
		Role:           domain.UserRole(sqlcUser.Role),
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}
