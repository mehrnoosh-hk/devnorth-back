package middleware

import (
	"net/http"
	"time"
)

// Timeout creates a middleware that sets a timeout for request processing
// This timeout applies to the entire request lifecycle (from handler start to response)
// If the handler exceeds the timeout, http.TimeoutHandler will:
//   - Send a 503 Service Unavailable response with the timeout message
//   - Cancel the handler's context (context-aware operations will stop)
//   - Prevent any subsequent writes to the response writer
func Timeout(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, timeout, "Request timeout")
	}
}
