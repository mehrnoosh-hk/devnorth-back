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
func NewRouter(userUseCase domain.UserUseCase, logger *slog.Logger, handlerTimeout time.Duration) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Timeout(handlerTimeout)) // Apply timeout first
	r.Use(middleware.Logger(logger))
	r.Use(middleware.CORS())

	// Initialize response writer
	responseWriter := response.NewWriter(logger)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(userUseCase, logger, responseWriter)
	healthHandler := handler.NewHealthHandler(responseWriter)

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

	return r
}
