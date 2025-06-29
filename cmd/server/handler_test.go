package server_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/victorvcruz/shipment-coordinator/cmd/server"
	"github.com/victorvcruz/shipment-coordinator/internal/carrier"
	carriermock "github.com/victorvcruz/shipment-coordinator/internal/carrier/mocks"
	"github.com/victorvcruz/shipment-coordinator/internal/order"
	ordermock "github.com/victorvcruz/shipment-coordinator/internal/order/mocks"
	"github.com/victorvcruz/shipment-coordinator/internal/shipping"
	shippingmock "github.com/victorvcruz/shipment-coordinator/internal/shipping/mocks"
	"github.com/victorvcruz/shipment-coordinator/pkg/states"
)

func TestHandler_CreateOrder_Success(t *testing.T) {
	orderSvc := &ordermock.ServiceMock{
		CreateFunc: func(ctx context.Context, o *order.Order) (*order.Order, error) {
			o.ID = uuid.New()
			return o, nil
		},
	}
	h := server.NewHandler(orderSvc, nil, nil)
	input := &order.CreateOrderInput{
		Body: order.CreateOrderInputBody{
			Product:       "Test",
			WeightKg:      2.5,
			DestinationUF: "SP",
		},
	}
	resp, err := h.CreateOrder(context.Background(), input)
	assert.NoError(t, err)
	assert.Equal(t, "Test", resp.Body.Product)
	assert.Equal(t, "SP", resp.Body.DestinationUF)
	assert.Equal(t, order.StatusCreated, resp.Body.Status)
}

func TestHandler_CreateOrder_InvalidUF(t *testing.T) {
	h := server.NewHandler(nil, nil, nil)
	input := &order.CreateOrderInput{
		Body: order.CreateOrderInputBody{
			Product:       "Test",
			WeightKg:      2.5,
			DestinationUF: "XX",
		},
	}
	resp, err := h.CreateOrder(context.Background(), input)
	assert.Nil(t, resp)
	assert.NotNil(t, err)
}

func TestHandler_GetOrder_Success(t *testing.T) {
	id := uuid.New()
	orderSvc := &ordermock.ServiceMock{
		GetByIDFunc: func(ctx context.Context, oid uuid.UUID) (*order.Order, error) {
			return &order.Order{
				ID:            id,
				Product:       "Test",
				Status:        order.StatusCreated,
				DestinationUF: states.SP,
				WeightKg:      decimal.NewFromFloat(1),
			}, nil
		},
	}
	h := server.NewHandler(orderSvc, nil, nil)
	input := &order.GetOrderParams{ID: id}
	resp, err := h.GetOrder(context.Background(), input)
	assert.NoError(t, err)
	assert.Equal(t, id, resp.Body.ID)
}

func TestHandler_GetOrder_NotFound(t *testing.T) {
	orderSvc := &ordermock.ServiceMock{
		GetByIDFunc: func(ctx context.Context, oid uuid.UUID) (*order.Order, error) {
			return nil, order.ErrOrderNotFound
		},
	}
	h := server.NewHandler(orderSvc, nil, nil)
	input := &order.GetOrderParams{ID: uuid.New()}
	resp, err := h.GetOrder(context.Background(), input)
	assert.Nil(t, resp)
	assert.NotNil(t, err)
}

func TestHandler_UpdateOrderStatus_Success(t *testing.T) {
	id := uuid.New()
	orderSvc := &ordermock.ServiceMock{
		UpdateStatusFunc: func(ctx context.Context, oid uuid.UUID, status order.Status) error {
			return nil
		},
	}
	h := server.NewHandler(orderSvc, nil, nil)
	input := &order.UpdateOrderStatusInput{
		ID: id,
		Body: order.UpdateOrderStatusInputBody{
			Status: string(order.StatusAwaitingPickup),
		},
	}
	resp, err := h.UpdateOrderStatus(context.Background(), input)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.Status)
}

func TestHandler_UpdateOrderStatus_InvalidStatus(t *testing.T) {
	h := server.NewHandler(nil, nil, nil)
	input := &order.UpdateOrderStatusInput{
		ID: uuid.New(),
		Body: order.UpdateOrderStatusInputBody{
			Status: "invalid",
		},
	}
	resp, err := h.UpdateOrderStatus(context.Background(), input)
	assert.Nil(t, resp)
	assert.NotNil(t, err)
}

func TestHandler_UpdateOrderStatus_AlreadySet(t *testing.T) {
	orderSvc := &ordermock.ServiceMock{
		UpdateStatusFunc: func(ctx context.Context, oid uuid.UUID, status order.Status) error {
			return order.ErrStatusAlreadySet
		},
	}
	h := server.NewHandler(orderSvc, nil, nil)
	input := &order.UpdateOrderStatusInput{
		ID: uuid.New(),
		Body: order.UpdateOrderStatusInputBody{
			Status: string(order.StatusCreated),
		},
	}
	resp, err := h.UpdateOrderStatus(context.Background(), input)
	assert.Nil(t, resp)
	assert.NotNil(t, err)
}

func TestHandler_CreateCarrier_Success(t *testing.T) {
	carrierSvc := &carriermock.ServiceMock{
		CreateFunc: func(ctx context.Context, c *carrier.Carrier) (*carrier.Carrier, error) {
			c.ID = uuid.New()
			return c, nil
		},
	}
	h := server.NewHandler(nil, carrierSvc, nil)
	input := &carrier.CreateCarrierInput{
		Body: carrier.CreateCarrierInputBody{
			Name: "Carrier1",
			Policies: []carrier.CarrierPolicyInput{
				{Region: "Sudeste", EstimatedDays: 2, PricePerKg: 10.00},
			},
		},
	}
	resp, err := h.CreateCarrier(context.Background(), input)
	assert.NoError(t, err)
	assert.Equal(t, "Carrier1", resp.Body.Name)
	assert.Equal(t, 201, resp.Status)
}

func TestHandler_CreateCarrier_InvalidRegion(t *testing.T) {
	h := server.NewHandler(nil, &carriermock.ServiceMock{}, nil)
	input := &carrier.CreateCarrierInput{
		Body: carrier.CreateCarrierInputBody{
			Name: "Carrier1",
			Policies: []carrier.CarrierPolicyInput{
				{Region: "Invalid", EstimatedDays: 2, PricePerKg: 10.00},
			},
		},
	}
	resp, err := h.CreateCarrier(context.Background(), input)
	assert.Nil(t, resp)
	assert.NotNil(t, err)
}

func TestHandler_GetQuotes_Success(t *testing.T) {
	orderID := uuid.New()
	shippingSvc := &shippingmock.ServiceMock{
		QuoteAllFunc: func(ctx context.Context, id uuid.UUID) ([]*shipping.Quote, error) {
			return []*shipping.Quote{
				{
					CarrierID:     uuid.New(),
					CarrierName:   "CarrierX",
					Price:         decimal.NewFromFloat(20),
					EstimatedDays: 3,
				},
			}, nil
		},
	}
	h := server.NewHandler(nil, nil, shippingSvc)
	input := &shipping.GetQuotesInput{OrderID: orderID.String()}
	resp, err := h.GetQuotes(context.Background(), input)
	assert.NoError(t, err)
	assert.Len(t, resp.Body, 1)
	assert.Equal(t, "CarrierX", resp.Body[0].CarrierName)
}

func TestHandler_GetQuotes_InvalidID(t *testing.T) {
	h := server.NewHandler(nil, nil, &shippingmock.ServiceMock{})
	input := &shipping.GetQuotesInput{OrderID: "invalid-uuid"}
	resp, err := h.GetQuotes(context.Background(), input)
	assert.Nil(t, resp)
	assert.NotNil(t, err)
}

func TestHandler_ContractCarrier_Success(t *testing.T) {
	shippingSvc := &shippingmock.ServiceMock{
		ContractCarrierFunc: func(ctx context.Context, orderID, carrierID uuid.UUID) (*shipping.Contract, error) {
			return &shipping.Contract{
				ID:            uuid.New(),
				OrderID:       orderID,
				CarrierID:     carrierID,
				Price:         decimal.NewFromFloat(100),
				EstimatedDays: 2,
				ContractedAt:  time.Now().UTC(),
				CreatedAt:     time.Now().UTC(),
				UpdatedAt:     time.Now().UTC(),
			}, nil
		},
	}
	h := server.NewHandler(nil, nil, shippingSvc)
	input := &shipping.ContractCarrierInput{
		Body: shipping.ContractCarrierInputBody{
			OrderID:   uuid.New(),
			CarrierID: uuid.New(),
		},
	}
	resp, err := h.ContractCarrier(context.Background(), input)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.Status)
}

func TestHandler_ContractCarrier_NoValidPolicy(t *testing.T) {
	shippingSvc := &shippingmock.ServiceMock{
		ContractCarrierFunc: func(ctx context.Context, orderID, carrierID uuid.UUID) (*shipping.Contract, error) {
			return nil, shipping.ErrNoValidPolicy
		},
	}
	h := server.NewHandler(nil, nil, shippingSvc)
	input := &shipping.ContractCarrierInput{
		Body: shipping.ContractCarrierInputBody{
			OrderID:   uuid.New(),
			CarrierID: uuid.New(),
		},
	}
	resp, err := h.ContractCarrier(context.Background(), input)
	assert.Nil(t, resp)
	assert.NotNil(t, err)
}
