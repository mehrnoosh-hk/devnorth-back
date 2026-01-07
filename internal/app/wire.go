package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mehrnoosh-hk/devnorth-back/config"
	"github.com/mehrnoosh-hk/devnorth-back/internal/database"
	httpDelivery "github.com/mehrnoosh-hk/devnorth-back/internal/delivery/http"
	"github.com/mehrnoosh-hk/devnorth-back/internal/domain"
	"github.com/mehrnoosh-hk/devnorth-back/internal/security"
	"github.com/mehrnoosh-hk/devnorth-back/internal/usecase"
)

// initLogger initializes the structured logger
func initLogger(env string) *slog.Logger {
	var handler slog.Handler

	if strings.EqualFold(env, "production") {
		// JSON handler for production
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		// Text handler for development
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	return slog.New(handler)
}

// initDatabase initializes the database connection
func initDatabase(ctx context.Context, cfg config.DatabaseConfig, logger *slog.Logger) (*pgxpool.Pool, error) {
	db, err := database.NewConnection(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInitDB, err)
	}
	logger.Info("Database connection established")
	return db, nil
}

// initSecurity initializes security-related dependencies
func initSecurity(cfg config.JWTConfig, logger *slog.Logger) (domain.PasswordHasher, domain.TokenGenerator, error) {
	passwordHasher, err := security.NewBcryptHasher(10) // Cost factor 10
	if err != nil {
		logger.Error("Failed to wire dependency: password hasher", "Error", err)
		return nil, nil, fmt.Errorf("%w: %w", ErrInitPasswordHasher, err)
	}

	// TODO: This is a hardcoded secret key "key1", acceptable for POC
	// For production use acceptable key-rotation algorithm
	tokenGenerator, err := security.NewJWTGenerator("key1", cfg.Keys, cfg.TokenDuration, logger)
	if err != nil {
		logger.Error("Failed to wire dependency: token generator", "Error", err)
		return nil, nil, fmt.Errorf("%w: %w", ErrInitTokenGenerator, err)
	}

	return passwordHasher, tokenGenerator, nil
}

// initUseCases initializes application use cases
func initUseCases(
	userRepo domain.UserRepository,
	competencyRepo domain.CompetencyRepository,
	passwordHasher domain.PasswordHasher,
	tokenGenerator domain.TokenGenerator,
	logger *slog.Logger,
) (domain.UserUseCase, domain.CompetencyUseCase, error) {
	userUseCase, err := usecase.NewUserUseCase(userRepo, passwordHasher, tokenGenerator, logger)
	if err != nil {
		logger.Error("Failed to wire dependency: user use case", "Error", err)
		return nil, nil, fmt.Errorf("%w: %w", ErrInitUserUseCase, err)
	}
	logger.Info("User use case initialized")

	competencyUseCase, err := usecase.NewCompetencyUseCase(competencyRepo, logger)
	if err != nil {
		logger.Error("Failed to wire dependency: competency use case", "Error", err)
		return nil, nil, fmt.Errorf("%w: %w", ErrInitCompetencyUseCase, err)
	}
	logger.Info("Competency use case initialized")

	return userUseCase, competencyUseCase, nil
}

// initServer initializes the HTTP server
func initServer(cfg config.ServerConfig, userUseCase domain.UserUseCase, competencyUseCase domain.CompetencyUseCase, logger *slog.Logger) (*httpDelivery.Server, error) {
	// Setup HTTP router with timeout from config
	handlerTimeout := time.Duration(cfg.HandlerTimeout) * time.Second
	router, err := httpDelivery.NewRouter(userUseCase, competencyUseCase, logger, handlerTimeout)
	if err != nil {
		return nil, err
	}

	// Create HTTP server
	return httpDelivery.NewServer(cfg, router, logger), nil
}