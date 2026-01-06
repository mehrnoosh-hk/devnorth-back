package ratelimit

import (
	"maps"
	"sync"
	"time"
)

// RateLimit defines the maximum number of requests allowed within a time window
type RateLimit struct {
	MaxRequests int           // Maximum number of requests allowed
	Window      time.Duration // Time window for rate limiting (sliding window)
}

// Limiter implements a sliding window rate limiter that tracks requests per IP and path
type Limiter struct {
	mu     sync.Mutex
	limits map[string]RateLimit // Per-path rate limit configuration
	// store maps IP addresses to their request history per path
	// Structure: IP -> Path -> []timestamp of requests
	store map[string]map[string][]time.Time
}

// Allow checks if a request from the given IP to the given path should be allowed
// Returns true if the request is within rate limits, false otherwise
func (l *Limiter) Allow(ip, path string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	// If path has no rate limit configured, allow the request
	rateLimit, exists := l.limits[path]
	if !exists {
		return true
	}

	// Initialize IP entry if this is the first request from this IP
	if l.store[ip] == nil {
		l.store[ip] = make(map[string][]time.Time)
	}

	// Get request history for this IP and path
	// Note: Returns nil if path not accessed before, which is safe in Go
	history := l.store[ip][path]

	// Filter out expired requests (outside the sliding window)
	validRequests := l.filterValidRequests(history, rateLimit.Window)
	l.store[ip][path] = validRequests

	// Check if adding this request would exceed the limit
	if len(validRequests) >= rateLimit.MaxRequests {
		return false
	}

	// Allow the request and record it
	l.store[ip][path] = append(validRequests, time.Now())
	return true
}

// filterValidRequests returns only the requests that fall within the current time window
func (l *Limiter) filterValidRequests(requests []time.Time, window time.Duration) []time.Time {
	if len(requests) == 0 {
		return []time.Time{}
	}

	now := time.Now()
	cutoff := now.Add(-window)

	valid := make([]time.Time, 0, len(requests))
	for _, timestamp := range requests {
		if timestamp.After(cutoff) {
			valid = append(valid, timestamp)
		}
	}

	return valid
}

// AddLimit adds or updates rate limits for multiple paths
// This method is thread-safe and can be called during runtime
func (l *Limiter) AddLimit(limitList map[string]RateLimit) {
	l.mu.Lock()
	defer l.mu.Unlock()
	maps.Copy(l.limits, limitList)
}

// NewLimiter creates a new rate limiter instance
// Use AddLimit to configure rate limits for specific paths
func NewLimiter() *Limiter {
	return &Limiter{
		limits: make(map[string]RateLimit),
		store:  make(map[string]map[string][]time.Time),
	}
}
