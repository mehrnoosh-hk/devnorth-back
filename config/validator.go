package config

import (
	"errors"
	"fmt"
	"strings"
)

// ErrConfigValidationFailed is returned when configuration validation fails
var ErrConfigValidationFailed = errors.New("configuration validation failed")

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if err := c.validateApp(); err != nil {
		return fmt.Errorf("%w: app config: %w", ErrConfigValidationFailed, err)
	}

	if err := c.validateServer(); err != nil {
		return fmt.Errorf("%w: server config: %w", ErrConfigValidationFailed, err)
	}

	if err := c.validateJWT(); err != nil {
		return fmt.Errorf("%w: JWT config: %w", ErrConfigValidationFailed, err)
	}

	if err := c.validateDatabase(); err != nil {
		return fmt.Errorf("%w: database config: %w", ErrConfigValidationFailed, err)
	}

	return nil
}

// validateServer validates server configuration
func (c *Config) validateServer() error {
	if c.Server.Port == "" {
		return errors.New("port is required")
	}

	if c.Server.ReadTimeout <= 0 {
		return errors.New("read timeout must be greater than 0")
	}

	if c.Server.WriteTimeout <= 0 {
		return errors.New("write timeout must be greater than 0")
	}

	if c.Server.HandlerTimeout <= 0 {
		return errors.New("handler timeout must be greater than 0")
	}

	return nil
}

// validateJWT validates JWT configuration
func (c *Config) validateJWT() error {
	if len(c.JWT.Keys) == 0 {
		return errors.New("at least one JWT key is required")
	}

	// Check if at least one key is non-empty
	hasValidKey := false
	for kid, key := range c.JWT.Keys {
		if strings.TrimSpace(key) != "" {
			hasValidKey = true
		} else if key == "" {
			return fmt.Errorf("JWT key '%s' is empty", kid)
		}
	}

	if !hasValidKey {
		return errors.New("at least one valid JWT key is required")
	}

	if c.JWT.TokenDuration <= 0 {
		return errors.New("token duration must be greater than 0")
	}

	return nil
}

// validateDatabase validates database configuration
func (c *Config) validateDatabase() error {
	// If URL is provided, we can skip individual field validation
	if c.Database.URL != "" {
		return c.validateDatabaseSettings()
	}

	// Otherwise validate individual connection fields
	if c.Database.Host == "" {
		return errors.New("database host is required")
	}

	if c.Database.User == "" {
		return errors.New("database user is required")
	}

	if c.Database.Name == "" {
		return errors.New("database name is required")
	}

	return c.validateDatabaseSettings()
}

// validateDatabaseSettings validates database pool and timeout settings
func (c *Config) validateDatabaseSettings() error {
	if c.Database.MaxOpenConns <= 0 {
		return errors.New("max open connections must be greater than 0")
	}

	if c.Database.MinOpenConns < 0 {
		return errors.New("min open connections cannot be negative")
	}

	if c.Database.MinOpenConns > c.Database.MaxOpenConns {
		return errors.New("min open connections cannot exceed max open connections")
	}

	if c.Database.ConnMaxLifetime < 0 {
		return errors.New("connection max lifetime cannot be negative")
	}

	if c.Database.ConnMaxIdleTime < 0 {
		return errors.New("connection max idle time cannot be negative")
	}

	if c.Database.ConnectTimeout <= 0 {
		return errors.New("connect timeout must be greater than 0")
	}

	if c.Database.QueryTimeout <= 0 {
		return errors.New("query timeout must be greater than 0")
	}

	return nil
}

// validateApp validates application-level configuration
func (c *Config) validateApp() error {
	if c.App.ShutdownTimeout <= 0 {
		return errors.New("shutdown timeout must be greater than 0")
	}

	return nil
}
