package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"time"

	"github.com/joho/godotenv"
	"github.com/mehrnoosh-hk/devnorth-back/config"
	"github.com/mehrnoosh-hk/devnorth-back/internal/app"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// TODO: move it to config validation
	// Validate JWT secret
	if cfg.JWT.SecretKey == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	// Initialize application context
	ctx := context.Background()

	// Initialize app (wire dependencies: database, logger, etc.)
	application, err := app.NewApp(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer application.Close()

	// Start server in goroutine
	go func() {
		if err := application.Server.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	application.Logger.Info("Application started successfully",
		"environment", cfg.AppEnv,
		"server_address", cfg.Server.Host+":"+cfg.Server.Port,
	)

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	application.Logger.Info("Shutting down application...")

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := application.Server.Shutdown(shutdownCtx); err != nil {
		application.Logger.Error("Error during server shutdown", "error", err)
	}

	application.Logger.Info("Application stopped")
}
