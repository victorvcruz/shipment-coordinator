package order

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/victorvcruz/shipment-coordinator/internal/platform/telemetry"
	log "go.uber.org/zap"
)

var (
	ErrStatusAlreadySet        = errors.New("status already set")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
)

//go:generate moq -pkg mocks -out mocks/service.go . Service
type Service interface {
	Create(ctx context.Context, order *Order) (*Order, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status Status) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, order *Order) (*Order, error) {
	id, err := s.repo.Create(ctx, order)
	if err != nil {
		log.L().
			Error("failed to create order", log.String("product", order.Product), log.Error(err))
		return nil, err
	}

	telemetry.OrderCreatedCounter.Add(ctx, 1)

	order.ID = id
	return order, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*Order, error) {
	found, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.L().
			Error("failed to get order by ID", log.String("order_id", id.String()), log.Error(err))
		return nil, err
	}
	return found, err
}

func (s *service) UpdateStatus(ctx context.Context, id uuid.UUID, status Status) error {
	order, err := s.GetByID(ctx, id)
	if err != nil {
		log.L().
			Error("failed to get order for status update", log.String("order_id", id.String()), log.Error(err))
		return err
	}

	if order.Status == status {
		log.L().
			Info("status is already set", log.String("order_id", id.String()), log.String("status", string(status)))
		return ErrStatusAlreadySet
	}

	if StatusOrder[status] < StatusOrder[order.Status] {
		log.L().
			Error("invalid status transition", log.String("order_id", id.String()), log.String("current_status", string(order.Status)), log.String("new_status", string(status)))
		return ErrInvalidStatusTransition
	}

	err = s.repo.UpdateStatus(ctx, id, status)
	if err != nil {
		log.L().
			Error("failed to update status", log.String("order_id", id.String()), log.Error(err))
		return err
	}
	return err
}
