package database

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"example.com/go-backend-foundation/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(
	ctx context.Context,
	cfg config.DatabaseConfig,
) (*pgxpool.Pool, error) {

	// Build PostgreSQL connection URL safely.
	connectionURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host: net.JoinHostPort(
			cfg.Host,
			strconv.Itoa(cfg.Port),
		),
		Path: "/" + cfg.Database,
	}

	query := connectionURL.Query()
	query.Set("sslmode", "disable")
	connectionURL.RawQuery = query.Encode()

	// Parse PostgreSQL connection configuration.
	poolConfig, err := pgxpool.ParseConfig(connectionURL.String())
	if err != nil {
		return nil, fmt.Errorf("parse database config: %w", err)
	}

	// Configure connection pool.
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = time.Minute

	// Create connection pool.
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}

	// Verify database connectivity.
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return pool, nil
}
