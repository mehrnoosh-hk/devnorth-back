package http

import (
	"log/slog"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mehrnoosh-hk/devnorth-back/internal/delivery/http/handler"
	"github.com/mehrnoosh-hk/devnorth-back/internal/delivery/http/middleware"
	"github.com/mehrnoosh-hk/devnorth-back/internal/delivery/http/response"
	"github.com/mehrnoosh-hk/devnorth-back/internal/domain"
)

// NewRouter creates and configures the HTTP router
func NewRouter(userUseCase domain.UserUseCase, logger *slog.Logger, handlerTimeout time.Duration) (*chi.Mux, error) {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Timeout(handlerTimeout)) // Apply timeout first
	r.Use(middleware.Logger(logger))
	r.Use(middleware.CORS())

	// Initialize response writer
	responseWriter, err := response.NewWriter(logger)
	if err != nil {
		return nil, err
	}

	// Initialize handlers
	authHandler, err := handler.NewAuthHandler(userUseCase, logger, responseWriter)
	if err != nil {
		return nil, err
	}
	healthHandler, err := handler.NewHealthHandler(responseWriter)
	if err != nil {
		return nil, err
	}

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Health check
		r.Get("/health", healthHandler.Check)

		// Authentication routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
		})
	})

	return r, nil
}
