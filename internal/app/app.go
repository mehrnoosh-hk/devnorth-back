package app

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mehrnoosh-hk/devnorth-back/config"
	"github.com/mehrnoosh-hk/devnorth-back/internal/database"
	httpDelivery "github.com/mehrnoosh-hk/devnorth-back/internal/delivery/http"
	"github.com/mehrnoosh-hk/devnorth-back/internal/domain"
	"github.com/mehrnoosh-hk/devnorth-back/internal/repository"
	"github.com/mehrnoosh-hk/devnorth-back/internal/security"
	"github.com/mehrnoosh-hk/devnorth-back/internal/usecase"
)

// App holds the application state and dependencies
type App struct {
	Config      *config.Config
	Logger      *slog.Logger
	DB          *pgxpool.Pool
	UserRepo    domain.UserRepository
	UserUseCase domain.UserUseCase
	Server      *httpDelivery.Server
}

// NewApp initializes and returns a new App instance with all dependencies
func NewApp(ctx context.Context, cfg *config.Config) (*App, error) {
	// Initialize logger
	logger := initLogger(cfg.AppEnv)

	// Initialize database connection
	db, err := database.NewConnection(ctx, cfg.Database)
	if err != nil {
		return nil, err
	}

	logger.Info("Database connection established")

	// TODO: Refactor to a Repository interface and its own initialization function when project grows
	userRepo := repository.NewUserRepository(db, logger)

	passwordHasher := security.NewBcryptHasher(10) // Cost factor 10
	tokenGenerator := security.NewJWTGenerator(cfg.JWT.SecretKey, cfg.JWT.TokenDuration)

	userUseCase := usecase.NewUserUseCase(userRepo, passwordHasher, tokenGenerator)

	// Setup HTTP router with timeout from config
	handlerTimeout := time.Duration(cfg.Server.HandlerTimeout) * time.Second
	router := httpDelivery.NewRouter(userUseCase, logger, handlerTimeout)

	// // Create HTTP server
	server := httpDelivery.NewServer(cfg.Server, router, logger)

	return &App{
		Config:      cfg,
		Logger:      logger,
		DB:          db,
		UserRepo:    userRepo,
		UserUseCase: userUseCase,
		Server:      server,
	}, nil
}

// Close gracefully shuts down the application
func (a *App) Close() {
	if a.DB != nil {
		if a.Logger != nil {
			a.Logger.Info("Closing database connection")
		}
		a.DB.Close() // pgxpool.Pool.Close() doesn't return an error
	}
}

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
