package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/victorvcruz/shipment-coordinator/internal/carrier"
	"github.com/victorvcruz/shipment-coordinator/internal/order"
	"github.com/victorvcruz/shipment-coordinator/internal/shipping"
	"github.com/victorvcruz/shipment-coordinator/pkg/states"
	log "go.uber.org/zap"
)

type Handler struct {
	orderService    order.Service
	carrierService  carrier.Service
	shippingService shipping.Service
}

func NewHandler(
	service order.Service,
	carrierService carrier.Service,
	shippingService shipping.Service,
) *Handler {
	return &Handler{
		orderService:    service,
		carrierService:  carrierService,
		shippingService: shippingService,
	}
}

func (h *Handler) CreateOrder(
	ctx context.Context,
	input *order.CreateOrderInput,
) (*order.OrderResponseOutput, error) {
	now := time.Now().UTC()
	state, ok := states.States[input.Body.DestinationUF]
	if !ok {
		return nil, huma.Error400BadRequest("invalid destination UF")
	}

	createdOrder, err := h.orderService.Create(ctx, &order.Order{
		Product:       input.Body.Product,
		WeightKg:      decimal.NewFromFloat(input.Body.WeightKg),
		DestinationUF: state,
		Status:        order.StatusCreated,
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to create order", err)
	}

	log.L().Info("Created order", log.String("order_id", createdOrder.ID.String()))
	return &order.OrderResponseOutput{
		Body: order.OrderResponseOutputBody{
			ID:            createdOrder.ID,
			Product:       createdOrder.Product,
			WeightKg:      createdOrder.WeightKg.StringFixed(2),
			DestinationUF: createdOrder.DestinationUF.Sigla,
			Status:        createdOrder.Status,
			CreatedAt:     createdOrder.CreatedAt,
			UpdatedAt:     createdOrder.UpdatedAt,
		},
		Status: http.StatusCreated,
	}, nil
}

func (h *Handler) GetOrder(
	ctx context.Context,
	input *order.GetOrderParams,
) (*order.OrderResponseOutput, error) {
	found, err := h.orderService.GetByID(ctx, input.ID)
	if err != nil {
		if errors.Is(err, order.ErrOrderNotFound) {
			return nil, huma.Error404NotFound("order not found")
		}
	}

	return &order.OrderResponseOutput{
		Body: order.OrderResponseOutputBody{
			ID:            found.ID,
			Product:       found.Product,
			WeightKg:      found.WeightKg.StringFixed(2),
			DestinationUF: found.DestinationUF.Sigla,
			Status:        found.Status,
			CreatedAt:     found.CreatedAt,
			UpdatedAt:     found.UpdatedAt,
		},
		Status: http.StatusOK,
	}, nil
}

func (h *Handler) UpdateOrderStatus(
	ctx context.Context,
	input *order.UpdateOrderStatusInput,
) (*order.UpdateOrderStatusOutput, error) {
	status, ok := order.StatusValues[input.Body.Status]
	if !ok {
		return nil, huma.Error400BadRequest("invalid status")
	}

	err := h.orderService.UpdateStatus(ctx, input.ID, status)
	if err != nil {
		switch {
		case errors.Is(err, order.ErrStatusAlreadySet):
			return nil, huma.Error400BadRequest("status already set")
		case errors.Is(err, order.ErrInvalidStatusTransition):
			return nil, huma.Error400BadRequest("invalid status transition")
		case errors.Is(err, order.ErrOrderNotFound):
			return nil, huma.Error404NotFound("order not found")
		}
		return nil, err
	}

	log.L().
		Info("Updated order status", log.String("order_id", input.ID.String()), log.String("status", input.Body.Status))
	return &order.UpdateOrderStatusOutput{
		Status: http.StatusOK,
	}, nil
}

func (h *Handler) CreateCarrier(
	ctx context.Context,
	input *carrier.CreateCarrierInput,
) (*carrier.CarrierResponseOutput, error) {
	now := time.Now().UTC()

	policies := make([]carrier.Policy, len(input.Body.Policies))
	for i, p := range input.Body.Policies {
		region, ok := states.Regions[p.Region]
		if !ok {
			return nil, huma.Error400BadRequest("invalid region")
		}

		policies[i] = carrier.Policy{
			Region:        region,
			EstimatedDays: p.EstimatedDays,
			PricePerKg:    decimal.NewFromFloat(p.PricePerKg),
			CreatedAt:     now,
			UpdatedAt:     now,
		}
	}

	createdCarrier, err := h.carrierService.Create(ctx, &carrier.Carrier{
		Name:      input.Body.Name,
		Policies:  policies,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to create order", err)
	}

	policiesResponse := make([]carrier.CarrierPolicyResponse, len(createdCarrier.Policies))
	for i, p := range createdCarrier.Policies {
		policiesResponse[i] = carrier.CarrierPolicyResponse{
			Region:        p.Region.Name,
			EstimatedDays: p.EstimatedDays,
			PricePerKg:    p.PricePerKg.StringFixed(2),
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
		}
	}

	log.L().Info("Created carrier", log.String("carrier_id", createdCarrier.ID.String()))
	return &carrier.CarrierResponseOutput{
		Body: carrier.CarrierResponseOutputBody{
			ID:        createdCarrier.ID,
			Name:      createdCarrier.Name,
			Policies:  policiesResponse,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Status: http.StatusCreated,
	}, nil
}

func (h *Handler) GetQuotes(
	ctx context.Context,
	input *shipping.GetQuotesInput,
) (*shipping.QuotesOutput, error) {
	id, err := uuid.Parse(input.OrderID)
	if err != nil {
		return nil, huma.Error404NotFound("invalid ID format")
	}

	quotes, err := h.shippingService.QuoteAll(ctx, id)
	if err != nil {
		if errors.Is(err, order.ErrOrderNotFound) {
			return nil, huma.Error404NotFound("order not found")
		}
		return nil, err
	}

	quotesResponse := make([]shipping.QuotesOutputBody, len(quotes))
	for i, q := range quotes {
		quotesResponse[i] = shipping.QuotesOutputBody{
			CarrierID:     q.CarrierID.String(),
			CarrierName:   q.CarrierName,
			Price:         q.Price.StringFixed(2),
			EstimatedDays: q.EstimatedDays,
		}
	}

	return &shipping.QuotesOutput{
		Body:   quotesResponse,
		Status: http.StatusOK,
	}, nil
}

func (h *Handler) ContractCarrier(
	ctx context.Context,
	input *shipping.ContractCarrierInput,
) (*shipping.ContractCarrierOutput, error) {
	contract, err := h.shippingService.ContractCarrier(
		ctx,
		input.Body.OrderID,
		input.Body.CarrierID,
	)
	if err != nil {
		switch {
		case errors.Is(err, order.ErrOrderNotFound):
			return nil, huma.Error404NotFound("order not found")
		case errors.Is(err, carrier.ErrCarrierNotFound):
			return nil, huma.Error404NotFound("carrier not found")
		case errors.Is(err, order.ErrInvalidStatusTransition):
			return nil, huma.Error400BadRequest("invalid order status for contracting")
		case errors.Is(err, shipping.ErrNoValidPolicy):
			return nil, huma.Error400BadRequest(
				"no valid policy for the carrier in the order's destination region",
			)
		}
		return nil, err
	}

	log.L().Info("Carrier contracted", log.String("contract_id", contract.ID.String()))
	return &shipping.ContractCarrierOutput{
		Body: shipping.ContractCarrierOutputBody{
			ID:            contract.ID.String(),
			OrderID:       contract.OrderID.String(),
			CarrierID:     contract.CarrierID.String(),
			Price:         contract.Price.StringFixed(2),
			EstimatedDays: contract.EstimatedDays,
			ContractedAt:  contract.ContractedAt,
			CreatedAt:     contract.CreatedAt,
			UpdatedAt:     contract.UpdatedAt,
		},
		Status: http.StatusOK,
	}, nil
}
