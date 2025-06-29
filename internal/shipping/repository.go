package shipping

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const queryInsertContract = `
	INSERT INTO contracts (order_id, carrier_id, price, estimated_days, contracted_at, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id
`

//go:generate moq -pkg mocks -out mocks/repository.go . Repository
type Repository interface {
	Insert(ctx context.Context, c *Contract) (uuid.UUID, error)
}

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &repository{pool: pool}
}

func (r *repository) Insert(ctx context.Context, c *Contract) (uuid.UUID, error) {
	var id uuid.UUID

	err := r.pool.QueryRow(ctx, queryInsertContract,
		c.OrderID,
		c.CarrierID,
		c.Price,
		c.EstimatedDays,
		c.ContractedAt,
		c.CreatedAt,
		c.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert contract: %w", err)
	}

	return id, nil
}
