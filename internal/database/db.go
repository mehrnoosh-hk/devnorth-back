package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mehrnoosh-hk/devnorth-back/config"
)

// NewConnection initializes a pgxpool connection pool with configuration settings
func NewConnection(ctx context.Context, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	// Parse and configure connection
	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Configure connection pool settings
	poolConfig.MaxConns = int32(cfg.MaxOpenConns)
	poolConfig.MinConns = int32(cfg.MinOpenConns)
	poolConfig.MaxConnLifetime = time.Duration(cfg.ConnMaxLifetime) * time.Minute
	poolConfig.MaxConnIdleTime = time.Duration(cfg.ConnMaxIdleTime) * time.Minute
	poolConfig.HealthCheckPeriod = 1 * time.Minute
	poolConfig.ConnConfig.ConnectTimeout = time.Duration(cfg.ConnectTimeout) * time.Second

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Verify connection with timeout
	pingCtx, cancel := context.WithTimeout(ctx, time.Duration(cfg.ConnectTimeout)*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}
