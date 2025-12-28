package app

import (
	"log/slog"
	"os"

	"github.com/mehrnoosh-hk/devnorth-back/config"
)

// App holds the application state and dependencies
type App struct {
	Config *config.Config
	Logger *slog.Logger
}

// NewApp initializes and returns a new App instance with all dependencies
func NewApp(cfg *config.Config) (*App, error) {
	// Initialize logger
	logger := initLogger(cfg.AppEnv)

	return &App{
		Config: cfg,
		Logger: logger,
	}, nil
}

// Close gracefully shuts down the application
// func (a *App) Close() error {
// 	if a.DB != nil {
// 		if err := a.DB.Close(); err != nil {
// 			return fmt.Errorf("failed to close database: %w", err)
// 		}
// 	}
// 	return nil
// }

// initLogger initializes the structured logger
func initLogger(env string) *slog.Logger {
	var handler slog.Handler

	if env == "production" {
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