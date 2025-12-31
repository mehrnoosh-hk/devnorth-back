package config

import (
	"fmt"
	"os"
	"strconv"
)

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host string
	Port string
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	// Connection details
	Host     string
	User     string
	Password string
	Name     string
	SSLMode  string
	URL      string

	// Connection pool settings
	MaxOpenConns    int // Maximum number of open connections to the database
	MaxIdleConns    int // Maximum number of idle connections in the pool
	ConnMaxLifetime int // Maximum lifetime of a connection in minutes
	ConnMaxIdleTime int // Maximum time a connection can be idle in minutes

	// Timeout settings
	ConnectTimeout int // Connection timeout in seconds
	QueryTimeout   int // Query execution timeout in seconds

	// Migration settings
	MigrationsPath  string // Path to migration files
	MigrationsTable string // Table name for tracking migrations
}

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	AppEnv   string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			// Connection details
			Host:     getEnv("DB_HOST", "localhost"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "devnorth"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			URL:      getEnv("DB_URL", ""),

			// Connection pool settings
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvAsInt("DB_CONN_MAX_LIFETIME", 5),
			ConnMaxIdleTime: getEnvAsInt("DB_CONN_MAX_IDLE_TIME", 5),

			// Timeout settings
			ConnectTimeout: getEnvAsInt("DB_CONNECT_TIMEOUT", 10),
			QueryTimeout:   getEnvAsInt("DB_QUERY_TIMEOUT", 30),

			// Migration settings
			MigrationsPath:  getEnv("DB_MIGRATIONS_PATH", "db/migrations"),
			MigrationsTable: getEnv("DB_MIGRATIONS_TABLE", "schema_migrations"),
		},
		AppEnv: getEnv("APP_ENV", "development"),
	}

	// Build DB_URL if not provided
	if cfg.Database.URL == "" {
		cfg.Database.URL = fmt.Sprintf(
			"postgresql://%s:%s@%s/%s?sslmode=%s",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Name,
			cfg.Database.SSLMode,
		)
	}

	return cfg, nil
}

// getEnv reads an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt reads an environment variable as integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
