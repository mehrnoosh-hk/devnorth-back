package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host           string
	Port           string
	ReadTimeout    int // Read timeout in seconds
	WriteTimeout   int // Write timeout in seconds
	HandlerTimeout int // Handler/request processing timeout in seconds
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
	MinOpenConns    int // Minimum number of open connections to the database
	ConnMaxLifetime int // Maximum lifetime of a connection in minutes
	ConnMaxIdleTime int // Maximum time a connection can be idle in minutes

	// Timeout settings
	ConnectTimeout int // Connection timeout in seconds
	QueryTimeout   int // Query execution timeout in seconds

	// Migration settings
	MigrationsPath  string // Path to migration files
	MigrationsTable string // Table name for tracking migrations
}

// JWTConfig holds JWT authentication configuration
type JWTConfig struct {
	Keys          map[string]string // Map of kid to key
	TokenDuration int               // Token expiration duration in minutes
}

// AppConfig holds application-level configuration
type AppConfig struct {
	Env             string // Application environment (development, staging, production)
	ShutdownTimeout int    // Graceful shutdown timeout in seconds
}

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	App      AppConfig
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Host:           getEnv("SERVER_HOST", "localhost"),
			Port:           getEnv("SERVER_PORT", "8080"),
			ReadTimeout:    getEnvAsInt("SERVER_READ_TIMEOUT", 15),
			WriteTimeout:   getEnvAsInt("SERVER_WRITE_TIMEOUT", 15),
			HandlerTimeout: getEnvAsInt("SERVER_HANDLER_TIMEOUT", 10),
		},
		JWT: JWTConfig{
			Keys: map[string]string{
				"key1": getEnv("JWT_KEY1", ""),
				"key2": getEnv("JWT_KEY2", ""),
			},
			TokenDuration: getEnvAsInt("JWT_TOKEN_DURATION", 15), // 15 minutes default for POC testing
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
			MinOpenConns:    getEnvAsInt("DB_MIN_OPEN_CONNS", 5),
			ConnMaxLifetime: getEnvAsInt("DB_CONN_MAX_LIFETIME", 5),
			ConnMaxIdleTime: getEnvAsInt("DB_CONN_MAX_IDLE_TIME", 5),

			// Timeout settings
			ConnectTimeout: getEnvAsInt("DB_CONNECT_TIMEOUT", 10),
			QueryTimeout:   getEnvAsInt("DB_QUERY_TIMEOUT", 30),

			// Migration settings
			MigrationsPath:  getEnv("DB_MIGRATIONS_PATH", "db/migrations"),
			MigrationsTable: getEnv("DB_MIGRATIONS_TABLE", "schema_migrations"),
		},
		App: AppConfig{
			Env:             getEnv("APP_ENV", "development"),
			ShutdownTimeout: getEnvAsInt("APP_SHUTDOWN_TIMEOUT", 30),
		},
	}

	// Build DB_URL if not provided
	if cfg.Database.URL == "" {
		u := &url.URL{
			Scheme: "postgres",
			User:   url.UserPassword(cfg.Database.User, cfg.Database.Password),
			Host:   cfg.Database.Host,
			Path:   cfg.Database.Name,
			RawQuery: url.Values{
				"sslmode": {cfg.Database.SSLMode},
			}.Encode(),
		}
		cfg.Database.URL = u.String()
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
		intValue, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("Warning: Invalid integer value for %s: %v\n", key, err)
			return defaultValue
		}
		return intValue
	}
	log.Printf("Warning: The value for %s is empty, using default %d", key, defaultValue)
	return defaultValue
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
		return nil, err
	}

	// Load configuration
	cfg, err := Load()
	if err != nil {
		log.Printf("Failed to load configuration: %v", err)
		return nil, err
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Printf("Invalid configuration: %v", err)
		return nil, err
	}
	return cfg, nil
}

func (c *Config) AppEnvironment() string {
	return c.App.Env
}

func (c *Config) ShutdownTimeout() time.Duration {
	return time.Duration(c.App.ShutdownTimeout) * time.Second
}

func (c *Config) ServerURL() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}
