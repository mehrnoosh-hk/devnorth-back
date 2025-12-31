package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mehrnoosh-hk/devnorth-back/config"
	"github.com/mehrnoosh-hk/devnorth-back/internal/app"
)

func main() {
	// Load app configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize application context with cancel
	ctxWithCancel, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize app (wire dependencies: database, logger, etc.)
	application, err := app.NewApp(ctxWithCancel, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	// Create error channel to receive server errors
	serverErr := make(chan error, 1)

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		err := application.Start()
		if err != nil {
			serverErr <- err
		}
	}()

	select {
	case err := <- serverErr:
		application.Logger.Error("Failed to start application server", "error", err)
	case sig := <- quit:
		application.Logger.Info("Received a shutdown signal", "Signal", sig)
	}
	
	shutdownContext, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout())
	defer shutdownCancel()

	cancel() // Signaling cancel to all components.
	application.Logger.Info("Shutting down application...")	
	application.Close(shutdownContext)
	application.Logger.Info("Application stopped")
}
