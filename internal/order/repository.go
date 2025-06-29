package order

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/victorvcruz/shipment-coordinator/internal/platform/telemetry"
	"github.com/victorvcruz/shipment-coordinator/pkg/states"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var ErrOrderNotFound = errors.New("order not found")

const (
	queryInsertOrder = `
		INSERT INTO orders (product, weight_kg, destination_uf, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	querySelectByID = `
		SELECT id, product, weight_kg, destination_uf, status, created_at, updated_at
		FROM orders
		WHERE id = $1
	`

	queryUpdateStatus = `
		UPDATE orders
		SET status = $2, updated_at = $3
		WHERE id = $1
		RETURNING id
	`
)

//go:generate moq -pkg mocks -out mocks/repository.go . Repository
type Repository interface {
	Create(ctx context.Context, order *Order) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status Status) error
}

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &repository{pool: pool}
}

func (r *repository) Create(ctx context.Context, order *Order) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var id uuid.UUID
	err := r.pool.QueryRow(ctx, queryInsertOrder,
		order.Product,
		order.WeightKg,
		order.DestinationUF.Sigla,
		order.Status,
		order.CreatedAt,
		order.UpdatedAt,
	).Scan(&id)

	return id, err
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := r.pool.QueryRow(ctx, querySelectByID, id)

	var (
		order Order
		state string
	)
	if err := row.Scan(
		&order.ID,
		&order.Product,
		&order.WeightKg,
		&state,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	); err != nil {
		return nil, ErrOrderNotFound
	}

	order.DestinationUF = states.States[state]

	return &order, nil
}

func (r *repository) UpdateStatus(ctx context.Context, id uuid.UUID, status Status) error {
	telemetry.OrderUpdatedCounter.Add(
		ctx,
		1,
		metric.WithAttributes(attribute.String("status", string(status))),
	)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var returnedID uuid.UUID
	err := r.pool.QueryRow(ctx, queryUpdateStatus, id, status, time.Now().UTC()).Scan(&returnedID)
	if err != nil {
		return ErrOrderNotFound
	}
	return nil
}
