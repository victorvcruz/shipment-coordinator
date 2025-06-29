package shipping

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/victorvcruz/shipment-coordinator/internal/carrier"
	"github.com/victorvcruz/shipment-coordinator/internal/order"
	"github.com/victorvcruz/shipment-coordinator/internal/platform/telemetry"
	log "go.uber.org/zap"
)

var ErrNoValidPolicy = errors.New(
	"no valid policy found for the carrier in the order's destination region",
)

//go:generate moq -pkg mocks -out mocks/service.go . Service
type Service interface {
	QuoteAll(ctx context.Context, orderID uuid.UUID) ([]*Quote, error)
	ContractCarrier(ctx context.Context, orderID, carrierID uuid.UUID) (*Contract, error)
}

type service struct {
	orderRepository    order.Repository
	carrierRepository  carrier.Repository
	shippingRepository Repository
}

func NewService(
	orderRepository order.Repository,
	carrierRepository carrier.Repository,
	shippingRepository Repository,
) Service {
	return &service{
		orderRepository:    orderRepository,
		carrierRepository:  carrierRepository,
		shippingRepository: shippingRepository,
	}
}

func (s *service) QuoteAll(ctx context.Context, orderID uuid.UUID) ([]*Quote, error) {
	order, err := s.orderRepository.GetByID(ctx, orderID)
	if err != nil {
		log.L().
			Error("failed to get order by ID", log.String("order_id", orderID.String()), log.Error(err))
		return nil, err
	}

	carriers, err := s.carrierRepository.ListAllByRegion(ctx, order.DestinationUF.Region)
	if err != nil {
		log.L().
			Error("failed to list carriers by region", log.String("region", order.DestinationUF.Region), log.Error(err))
		return nil, err
	}

	var quotes []*Quote
	for _, c := range carriers {

		var validPolicy carrier.Policy
		for _, policy := range c.Policies {
			if policy.Region.Name == order.DestinationUF.Region {
				validPolicy = policy
				break
			}
		}

		quote := &Quote{
			CarrierID:     c.ID,
			CarrierName:   c.Name,
			Price:         validPolicy.PricePerKg.Mul(order.WeightKg),
			EstimatedDays: validPolicy.EstimatedDays,
		}
		quotes = append(quotes, quote)
	}

	return quotes, nil
}

func (s *service) ContractCarrier(
	ctx context.Context,
	orderID, carrierID uuid.UUID,
) (*Contract, error) {
	o, err := s.orderRepository.GetByID(ctx, orderID)
	if err != nil {
		log.L().
			Error("failed to get order by ID", log.String("order_id", orderID.String()), log.Error(err))
		return nil, err
	}

	if o.Status != order.StatusCreated {
		log.L().
			Error("invalid order status for contracting", log.String("order_id", o.ID.String()), log.String("current_status", string(o.Status)))
		return nil, order.ErrInvalidStatusTransition
	}

	c, err := s.carrierRepository.GetByID(ctx, carrierID)
	if err != nil {
		log.L().
			Error("failed to get carrier by ID", log.String("carrier_id", carrierID.String()), log.Error(err))
		return nil, err
	}

	var validPolicy carrier.Policy
	for _, policy := range c.Policies {
		if policy.Region.Name == o.DestinationUF.Region {
			validPolicy = policy
			break
		}
	}

	if validPolicy.ID == uuid.Nil {
		return nil, ErrNoValidPolicy
	}

	now := time.Now().UTC()
	contract := &Contract{
		OrderID:       o.ID,
		CarrierID:     c.ID,
		Price:         validPolicy.PricePerKg.Mul(o.WeightKg),
		EstimatedDays: validPolicy.EstimatedDays,
		ContractedAt:  now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	id, err := s.shippingRepository.Insert(ctx, contract)
	if err != nil {
		log.L().
			Error("failed to insert contract", log.String("order_id", o.ID.String()), log.String("carrier_id", c.ID.String()), log.Error(err))
		return nil, err
	}
	contract.ID = id

	err = s.orderRepository.UpdateStatus(ctx, o.ID, order.StatusAwaitingPickup)
	if err != nil {
		log.L().
			Error("failed to get carrier by ID", log.String("carrier_id", carrierID.String()), log.Error(err))
		return nil, err
	}

	telemetry.ContractCreatedCounter.Add(ctx, 1)

	return contract, nil
}
