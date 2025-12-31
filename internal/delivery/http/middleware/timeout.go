package middleware

import (
	"context"
	"net/http"
	"time"
)

// Timeout creates a middleware that sets a timeout for request processing
// This timeout applies to the entire request lifecycle (from handler start to response)
func Timeout(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create a context with timeout
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			// Pass the request with timeout context to the next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
