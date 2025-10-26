package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect creates a connection pool with sane defaults.
func Connect(ctx context.Context, dsn string, maxConns int32) (*pgxpool.Pool, error) {
	if dsn == "" {
		return nil, fmt.Errorf("database dsn is required")
	}

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse dsn: %w", err)
	}

	if maxConns > 0 {
		cfg.MaxConns = maxConns
	}

	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	return pool, nil
}
