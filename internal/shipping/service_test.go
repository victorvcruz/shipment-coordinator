package shipping_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/victorvcruz/shipment-coordinator/internal/carrier"
	carriermock "github.com/victorvcruz/shipment-coordinator/internal/carrier/mocks"
	"github.com/victorvcruz/shipment-coordinator/internal/order"
	ordermock "github.com/victorvcruz/shipment-coordinator/internal/order/mocks"
	"github.com/victorvcruz/shipment-coordinator/internal/shipping"
	"github.com/victorvcruz/shipment-coordinator/internal/shipping/mocks"
	"github.com/victorvcruz/shipment-coordinator/pkg/states"
)

func TestService_QuoteAll_Success(t *testing.T) {
	ctx := context.Background()
	orderID := uuid.New()
	carrierID := uuid.New()

	state := states.SP
	region := states.Sudeste

	orderObj := &order.Order{
		ID:            orderID,
		WeightKg:      decimal.NewFromFloat(10),
		DestinationUF: state,
	}

	policy := carrier.Policy{
		ID:            uuid.New(),
		Region:        region,
		EstimatedDays: 3,
		PricePerKg:    decimal.NewFromFloat(5),
	}
	carrierObj := carrier.Carrier{
		ID:       carrierID,
		Name:     "CarrierX",
		Policies: []carrier.Policy{policy},
	}

	orderRepo := &ordermock.RepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*order.Order, error) {
			return orderObj, nil
		},
	}
	carrierRepo := &carriermock.RepositoryMock{
		ListAllByRegionFunc: func(ctx context.Context, region string) ([]carrier.Carrier, error) {
			return []carrier.Carrier{carrierObj}, nil
		},
	}
	shippingRepo := &mocks.RepositoryMock{}

	svc := shipping.NewService(orderRepo, carrierRepo, shippingRepo)
	quotes, err := svc.QuoteAll(ctx, orderID)
	assert.NoError(t, err)
	assert.Len(t, quotes, 1)
	assert.Equal(t, carrierID, quotes[0].CarrierID)
	assert.Equal(t, "CarrierX", quotes[0].CarrierName)
	assert.Equal(t, decimal.NewFromFloat(50), quotes[0].Price)
	assert.Equal(t, 3, quotes[0].EstimatedDays)
}

func TestService_QuoteAll_OrderNotFound(t *testing.T) {
	ctx := context.Background()
	orderRepo := &ordermock.RepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*order.Order, error) {
			return nil, errors.New("not found")
		},
	}
	carrierRepo := &carriermock.RepositoryMock{}
	shippingRepo := &mocks.RepositoryMock{}
	svc := shipping.NewService(orderRepo, carrierRepo, shippingRepo)
	quotes, err := svc.QuoteAll(ctx, uuid.New())
	assert.Error(t, err)
	assert.Nil(t, quotes)
}

func TestService_ContractCarrier_Success(t *testing.T) {
	ctx := context.Background()
	orderID := uuid.New()
	carrierID := uuid.New()

	state := states.SP
	region := states.Sudeste

	orderObj := &order.Order{
		ID:            orderID,
		WeightKg:      decimal.NewFromFloat(2),
		DestinationUF: state,
	}
	policy := carrier.Policy{
		ID:            uuid.New(),
		Region:        region,
		EstimatedDays: 5,
		PricePerKg:    decimal.NewFromFloat(7),
	}
	carrierObj := &carrier.Carrier{
		ID:       carrierID,
		Name:     "CarrierY",
		Policies: []carrier.Policy{policy},
	}
	contractID := uuid.New()
	orderRepo := &ordermock.RepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*order.Order, error) {
			return orderObj, nil
		},
	}
	carrierRepo := &carriermock.RepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*carrier.Carrier, error) {
			return carrierObj, nil
		},
	}
	shippingRepo := &mocks.RepositoryMock{
		InsertFunc: func(ctx context.Context, c *shipping.Contract) (uuid.UUID, error) {
			return contractID, nil
		},
	}
	svc := shipping.NewService(orderRepo, carrierRepo, shippingRepo)
	contract, err := svc.ContractCarrier(ctx, orderID, carrierID)
	assert.NoError(t, err)
	assert.Equal(t, contractID, contract.ID)
	assert.Equal(t, orderID, contract.OrderID)
	assert.Equal(t, carrierID, contract.CarrierID)
	assert.Equal(t, decimal.NewFromFloat(14), contract.Price)
	assert.Equal(t, 5, contract.EstimatedDays)
	assert.WithinDuration(t, time.Now().UTC(), contract.ContractedAt, time.Second)
}

func TestService_ContractCarrier_NoValidPolicy(t *testing.T) {
	ctx := context.Background()
	orderID := uuid.New()
	carrierID := uuid.New()

	state := states.SP

	orderObj := &order.Order{
		ID:            orderID,
		WeightKg:      decimal.NewFromFloat(1),
		DestinationUF: state,
	}
	carrierObj := &carrier.Carrier{
		ID:       carrierID,
		Name:     "CarrierZ",
		Policies: []carrier.Policy{},
	}
	orderRepo := &ordermock.RepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*order.Order, error) {
			return orderObj, nil
		},
	}
	carrierRepo := &carriermock.RepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*carrier.Carrier, error) {
			return carrierObj, nil
		},
	}
	shippingRepo := &mocks.RepositoryMock{}
	svc := shipping.NewService(orderRepo, carrierRepo, shippingRepo)
	contract, err := svc.ContractCarrier(ctx, orderID, carrierID)
	assert.ErrorIs(t, err, shipping.ErrNoValidPolicy)
	assert.Nil(t, contract)
}
