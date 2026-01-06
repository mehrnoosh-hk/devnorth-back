package app

import (
	"context"
	"log/slog"


	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mehrnoosh-hk/devnorth-back/config"
	httpDelivery "github.com/mehrnoosh-hk/devnorth-back/internal/delivery/http"
	"github.com/mehrnoosh-hk/devnorth-back/internal/domain"
	"github.com/mehrnoosh-hk/devnorth-back/internal/repository"
)

// App holds the application state and dependencies
type App struct {
	config      *config.Config
	Logger      *slog.Logger
	db          *pgxpool.Pool
	userRepo    domain.UserRepository
	userUseCase domain.UserUseCase
	server      *httpDelivery.Server
}

// NewApp initializes and returns a new App instance with all dependencies
func NewApp(ctx context.Context, cfg *config.Config) (*App, error) {
	// Initialize logger
	logger := initLogger(cfg.AppEnvironment())

	// Initialize database connection
	db, err := initDatabase(ctx, cfg.Database, logger)
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db, logger)

	// Initialize security dependencies
	passwordHasher, tokenGenerator, err := initSecurity(cfg.JWT, logger)
	if err != nil {
		db.Close()
		return nil, err
	}
	logger.Info("Security dependencies initialized")

	// Initialize use cases
	userUseCase, err := initUseCases(userRepo, passwordHasher, tokenGenerator, logger)
	if err != nil {
		db.Close()
		return nil, err
	}

	// Initialize HTTP server
	server, err := initServer(cfg.Server, userUseCase, logger)
	if err != nil {
		db.Close()
		return nil, err
	}

	return &App{
		config:      cfg,
		Logger:      logger,
		db:          db,
		userRepo:    userRepo,
		userUseCase: userUseCase,
		server:      server,
	}, nil
}

// Start encapsulate app server start
func (a *App) Start() error {
	return a.server.Start()
}

// Close gracefully shuts down the application
func (a *App) Close(ctx context.Context) {
	if a.db != nil {
		if a.Logger != nil {
			a.Logger.Info("Closing database connection")
		}
		a.db.Close() // pgxpool.Pool.Close() doesn't return an error
	}
	if a.server != nil {
		if a.Logger != nil {
			a.Logger.Info("Closing server")
		}
		a.server.Shutdown(ctx)
	}
}
