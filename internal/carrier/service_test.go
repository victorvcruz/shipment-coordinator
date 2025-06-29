package carrier_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/victorvcruz/shipment-coordinator/internal/carrier"
	"github.com/victorvcruz/shipment-coordinator/internal/carrier/mocks"
)

func TestService_Create_Success(t *testing.T) {
	ctx := context.Background()
	expectedID := uuid.New()
	c := &carrier.Carrier{Name: "Carrier Test"}

	repo := &mocks.RepositoryMock{
		CreateFunc: func(ctx context.Context, c *carrier.Carrier) (uuid.UUID, error) {
			return expectedID, nil
		},
	}

	svc := carrier.NewService(repo)
	created, err := svc.Create(ctx, c)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, created.ID)
	assert.Equal(t, "Carrier Test", created.Name)
}

func TestService_Create_Error(t *testing.T) {
	ctx := context.Background()
	c := &carrier.Carrier{Name: "Carrier Fail"}

	repo := &mocks.RepositoryMock{
		CreateFunc: func(ctx context.Context, c *carrier.Carrier) (uuid.UUID, error) {
			return uuid.Nil, errors.New("db error")
		},
	}

	svc := carrier.NewService(repo)
	created, err := svc.Create(ctx, c)
	assert.Error(t, err)
	assert.Nil(t, created)
}
