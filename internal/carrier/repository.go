package carrier

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/victorvcruz/shipment-coordinator/pkg/states"
)

var ErrCarrierNotFound = errors.New("carrier not found")

//go:generate moq -pkg mocks -out mocks/repository.go . Repository
type Repository interface {
	Create(ctx context.Context, carrier *Carrier) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Carrier, error)
	ListAll(ctx context.Context) ([]Carrier, error)
	ListAllByRegion(ctx context.Context, region string) ([]Carrier, error)
}

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{pool: db}
}

const (
	queryInsertCarrier = `
	INSERT INTO carriers (name, created_at, updated_at)
	VALUES ($1, $2, $3)
	RETURNING id
`
	queryGetCarrierByID = `
	SELECT c.id, c.name, c.created_at, c.updated_at,
		   p.id, p.carrier_id, p.region, p.estimated_days, p.price_per_kg, p.created_at, p.updated_at
	FROM carriers c
	LEFT JOIN carrier_policies p ON c.id = p.carrier_id
	WHERE c.id = $1
`

	queryInsertCarrierPolicy = `
	INSERT INTO carrier_policies (carrier_id, region, estimated_days, price_per_kg)
	VALUES ($1, $2, $3, $4)
`

	queryListCarriersWithPolicies = `
	SELECT c.id, c.name, c.created_at, c.updated_at,
		   p.region, p.estimated_days, p.price_per_kg
	FROM carriers c
	LEFT JOIN carrier_policies p ON c.id = p.carrier_id
`

	queryListCarriersByRegion = `
	SELECT c.id, c.name, c.created_at, c.updated_at,
		   p.region, p.estimated_days, p.price_per_kg
	FROM carriers c
	INNER JOIN carrier_policies p ON c.id = p.carrier_id
	WHERE p.region = $1
`
)

func (r *repository) Create(ctx context.Context, carrier *Carrier) (uuid.UUID, error) {
	now := time.Now()

	var carrierID uuid.UUID
	err := r.pool.QueryRow(ctx, queryInsertCarrier,
		carrier.Name,
		now,
		now,
	).Scan(&carrierID)
	if err != nil {
		return uuid.Nil, err
	}

	for _, p := range carrier.Policies {
		_, err := r.pool.Exec(ctx, queryInsertCarrierPolicy,
			carrierID,
			p.Region.Name,
			p.EstimatedDays,
			p.PricePerKg.String(),
		)
		if err != nil {
			return uuid.Nil, fmt.Errorf("error inserting policy for region %s: %w", p.Region, err)
		}
	}

	return carrierID, nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Carrier, error) {
	rows, err := r.pool.Query(ctx, queryGetCarrierByID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carrier *Carrier

	for rows.Next() {
		var (
			cID, policyID, policyCarrierID     uuid.UUID
			name                               string
			carrierCreatedAt, carrierUpdatedAt time.Time

			region                           *string
			estimatedDays                    *int
			priceStr                         *string
			policyCreatedAt, policyUpdatedAt *time.Time
		)

		if err := rows.Scan(
			&cID, &name, &carrierCreatedAt, &carrierUpdatedAt,
			&policyID, &policyCarrierID, &region, &estimatedDays, &priceStr, &policyCreatedAt, &policyUpdatedAt,
		); err != nil {
			return nil, err
		}

		if carrier == nil {
			carrier = &Carrier{
				ID:        cID,
				Name:      name,
				CreatedAt: carrierCreatedAt,
				UpdatedAt: carrierUpdatedAt,
				Policies:  []Policy{},
			}
		}

		if region != nil && estimatedDays != nil && priceStr != nil {
			priceDec, err := decimal.NewFromString(*priceStr)
			if err != nil {
				return nil, fmt.Errorf("invalid decimal price_per_kg: %w", err)
			}

			policy := Policy{
				ID:            policyID,
				CarrierID:     policyCarrierID,
				Region:        states.Regions[*region],
				EstimatedDays: *estimatedDays,
				PricePerKg:    priceDec,
			}

			if policyCreatedAt != nil {
				policy.CreatedAt = *policyCreatedAt
			}
			if policyUpdatedAt != nil {
				policy.UpdatedAt = *policyUpdatedAt
			}

			carrier.Policies = append(carrier.Policies, policy)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if carrier == nil {
		return nil, ErrCarrierNotFound
	}

	return carrier, nil
}

func (r *repository) ListAll(ctx context.Context) ([]Carrier, error) {
	rows, err := r.pool.Query(ctx, queryListCarriersWithPolicies)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	carrierMap := make(map[uuid.UUID]*Carrier)

	for rows.Next() {
		var (
			id            uuid.UUID
			name          string
			createdAt     time.Time
			updatedAt     time.Time
			region        *string
			estimatedDays *int
			priceStr      *string
		)

		err := rows.Scan(&id, &name, &createdAt, &updatedAt, &region, &estimatedDays, &priceStr)
		if err != nil {
			return nil, err
		}

		carrier, exists := carrierMap[id]
		if !exists {
			carrier = &Carrier{
				ID:        id,
				Name:      name,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
				Policies:  []Policy{},
			}
			carrierMap[id] = carrier
		}

		if region != nil && estimatedDays != nil && priceStr != nil {
			priceDec, err := decimal.NewFromString(*priceStr)
			if err != nil {
				return nil, fmt.Errorf("invalid decimal price_per_kg: %w", err)
			}

			policy := Policy{
				Region:        states.Regions[*region],
				EstimatedDays: *estimatedDays,
				PricePerKg:    priceDec,
			}
			carrier.Policies = append(carrier.Policies, policy)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	carriers := make([]Carrier, 0, len(carrierMap))
	for _, c := range carrierMap {
		carriers = append(carriers, *c)
	}

	return carriers, nil
}

func (r *repository) ListAllByRegion(ctx context.Context, region string) ([]Carrier, error) {
	rows, err := r.pool.Query(ctx, queryListCarriersByRegion, region)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	carrierMap := make(map[uuid.UUID]*Carrier)

	for rows.Next() {
		var (
			id            uuid.UUID
			name          string
			createdAt     time.Time
			updatedAt     time.Time
			regionStr     string
			estimatedDays int
			priceStr      string
		)

		err := rows.Scan(&id, &name, &createdAt, &updatedAt, &regionStr, &estimatedDays, &priceStr)
		if err != nil {
			return nil, err
		}

		carrier, exists := carrierMap[id]
		if !exists {
			carrier = &Carrier{
				ID:        id,
				Name:      name,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
				Policies:  []Policy{},
			}
			carrierMap[id] = carrier
		}

		priceDec, err := decimal.NewFromString(priceStr)
		if err != nil {
			return nil, fmt.Errorf("invalid decimal price_per_kg: %w", err)
		}

		policy := Policy{
			Region:        states.Regions[regionStr],
			EstimatedDays: estimatedDays,
			PricePerKg:    priceDec,
		}

		carrier.Policies = append(carrier.Policies, policy)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	carriers := make([]Carrier, 0, len(carrierMap))
	for _, c := range carrierMap {
		carriers = append(carriers, *c)
	}

	return carriers, nil
}
