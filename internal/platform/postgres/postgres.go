package postgres

import (
	"context"
	"fmt"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/victorvcruz/shipment-coordinator/internal/platform/config"
)

func Connect(ctx context.Context, config *config.AppConfig) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(config.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	cfg.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err = otelpgx.RecordStats(pool); err != nil {
		return nil, fmt.Errorf("record stats: %w", err)
	}

	return pool, nil
}
