package order_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/victorvcruz/shipment-coordinator/internal/order"
	"github.com/victorvcruz/shipment-coordinator/internal/order/mocks"
	"github.com/victorvcruz/shipment-coordinator/pkg/states"
)

func TestService_Create_Success(t *testing.T) {
	ctx := context.Background()
	expectedID := uuid.New()
	orderObj := &order.Order{
		Product:       "Test Product",
		WeightKg:      decimal.NewFromFloat(1.5),
		DestinationUF: states.SP,
		Status:        order.StatusCreated,
	}

	repo := &mocks.RepositoryMock{
		CreateFunc: func(ctx context.Context, o *order.Order) (uuid.UUID, error) {
			return expectedID, nil
		},
	}

	svc := order.NewService(repo)
	created, err := svc.Create(ctx, orderObj)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, created.ID)
	assert.Equal(t, "Test Product", created.Product)
}

func TestService_Create_Error(t *testing.T) {
	ctx := context.Background()
	orderObj := &order.Order{Product: "Fail"}
	repo := &mocks.RepositoryMock{
		CreateFunc: func(ctx context.Context, o *order.Order) (uuid.UUID, error) {
			return uuid.Nil, errors.New("db error")
		},
	}
	svc := order.NewService(repo)
	created, err := svc.Create(ctx, orderObj)
	assert.Error(t, err)
	assert.Nil(t, created)
}

func TestService_GetByID_Success(t *testing.T) {
	ctx := context.Background()
	expectedID := uuid.New()
	orderObj := &order.Order{ID: expectedID, Product: "Test"}
	repo := &mocks.RepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*order.Order, error) {
			return orderObj, nil
		},
	}
	svc := order.NewService(repo)
	found, err := svc.GetByID(ctx, expectedID)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, found.ID)
}

func TestService_GetByID_Error(t *testing.T) {
	ctx := context.Background()
	repo := &mocks.RepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*order.Order, error) {
			return nil, errors.New("not found")
		},
	}
	svc := order.NewService(repo)
	found, err := svc.GetByID(ctx, uuid.New())
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestService_UpdateStatus_Success(t *testing.T) {
	ctx := context.Background()
	orderID := uuid.New()
	orderObj := &order.Order{ID: orderID, Status: order.StatusCreated}
	repo := &mocks.RepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*order.Order, error) {
			return orderObj, nil
		},
		UpdateStatusFunc: func(ctx context.Context, id uuid.UUID, status order.Status) error {
			return nil
		},
	}
	svc := order.NewService(repo)
	err := svc.UpdateStatus(ctx, orderID, order.StatusAwaitingPickup)
	assert.NoError(t, err)
}

func TestService_UpdateStatus_AlreadySet(t *testing.T) {
	ctx := context.Background()
	orderID := uuid.New()
	orderObj := &order.Order{ID: orderID, Status: order.StatusCreated}
	repo := &mocks.RepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*order.Order, error) {
			return orderObj, nil
		},
	}
	svc := order.NewService(repo)
	err := svc.UpdateStatus(ctx, orderID, order.StatusCreated)
	assert.ErrorIs(t, err, order.ErrStatusAlreadySet)
}

func TestService_UpdateStatus_InvalidTransition(t *testing.T) {
	ctx := context.Background()
	orderID := uuid.New()
	orderObj := &order.Order{ID: orderID, Status: order.StatusShipped}
	repo := &mocks.RepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*order.Order, error) {
			return orderObj, nil
		},
	}
	svc := order.NewService(repo)
	err := svc.UpdateStatus(ctx, orderID, order.StatusCreated)
	assert.ErrorIs(t, err, order.ErrInvalidStatusTransition)
}

func TestService_UpdateStatus_RepoError(t *testing.T) {
	ctx := context.Background()
	orderID := uuid.New()
	orderObj := &order.Order{ID: orderID, Status: order.StatusCreated}
	repo := &mocks.RepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*order.Order, error) {
			return orderObj, nil
		},
		UpdateStatusFunc: func(ctx context.Context, id uuid.UUID, status order.Status) error {
			return errors.New("update error")
		},
	}
	svc := order.NewService(repo)
	err := svc.UpdateStatus(ctx, orderID, order.StatusAwaitingPickup)
	assert.Error(t, err)
}
