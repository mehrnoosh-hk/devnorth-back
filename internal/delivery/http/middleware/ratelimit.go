package middleware

import (
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/mehrnoosh-hk/devnorth-back/config"
)

// client tracks the rate limiter and last seen time for a single IP
type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter manages per-IP rate limiters
type RateLimiter struct {
	mu      sync.RWMutex
	clients map[string]*client
	rps     rate.Limit
	burst   int
	logger  *slog.Logger
}

// NewRateLimiter creates a RateLimiter and starts a background cleanup goroutine
func NewRateLimiter(cfg config.RateLimitConfig, logger *slog.Logger) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*client),
		rps:     rate.Limit(cfg.RPS),
		burst:   cfg.Burst,
		logger:  logger,
	}

	cleanupInterval := time.Duration(cfg.CleanupPeriod) * time.Second
	go rl.cleanup(cleanupInterval)

	return rl
}

// getLimiter retrieves or creates a rate limiter for the given IP
func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.RLock()
	c, exists := rl.clients[ip]
	rl.mu.RUnlock()

	if exists {
		rl.mu.Lock()
		c.lastSeen = time.Now()
		rl.mu.Unlock()
		return c.limiter
	}

	limiter := rate.NewLimiter(rl.rps, rl.burst)
	rl.mu.Lock()
	rl.clients[ip] = &client{limiter: limiter, lastSeen: time.Now()}
	rl.mu.Unlock()

	return limiter
}

// cleanup periodically removes stale client entries
func (rl *RateLimiter) cleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for ip, c := range rl.clients {
			if time.Since(c.lastSeen) > interval {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// rateLimitError mirrors dto.ErrorResponse to avoid a circular import
type rateLimitError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// RateLimit returns chi-compatible middleware that enforces per-IP rate limiting
func RateLimit(cfg config.RateLimitConfig, logger *slog.Logger) func(http.Handler) http.Handler {
	rl := NewRateLimiter(cfg, logger)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := clientIP(r)
			if !rl.getLimiter(ip).Allow() {
				logger.Warn("Rate limit exceeded", "ip", ip, "path", r.URL.Path)

				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", "1")
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(rateLimitError{
					Error:   "rate_limit_exceeded",
					Message: "Too many requests, please try again later",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// clientIP extracts the client IP from the request, checking X-Forwarded-For
// and X-Real-IP headers before falling back to RemoteAddr
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP in the comma-separated chain (original client)
		if i := strings.IndexByte(xff, ','); i > 0 {
			return strings.TrimSpace(xff[:i])
		}
		return strings.TrimSpace(xff)
	}

	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
