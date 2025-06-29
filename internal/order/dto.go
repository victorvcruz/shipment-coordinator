package order

import (
	"time"

	"github.com/google/uuid"
)

type CreateOrderInput struct {
	Body CreateOrderInputBody
}
type CreateOrderInputBody struct {
	Product       string  `json:"product"        required:"true" doc:"Product name"   example:"MacBook Pro 16"`
	WeightKg      float64 `json:"weight_kg"      required:"true" doc:"Weight kg"      example:"2.5"`
	DestinationUF string  `json:"destination_uf" required:"true" doc:"Destination UF" example:"SP"`
}

type OrderResponseOutput struct {
	Status int
	Body   OrderResponseOutputBody
}

type OrderResponseOutputBody struct {
	ID            uuid.UUID `json:"id"             doc:"Order ID"               example:"123e4567-e89b-12d3-a456-426614174000"`
	Product       string    `json:"product"        doc:"Product name"           example:"MacBook Pro 16"`
	WeightKg      string    `json:"weight_kg"      doc:"Weight in kg"           example:"2.5"`
	DestinationUF string    `json:"destination_uf" doc:"Destination UF"         example:"SP"`
	Status        Status    `json:"status"         doc:"Order status"           example:"created"`
	CreatedAt     time.Time `json:"created_at"     doc:"Order creation date"    example:"2023-10-01T12:00:00Z"`
	UpdatedAt     time.Time `json:"updated_at"     doc:"Order last update date" example:"2023-10-01T12:00:00Z"`
}

type GetOrderParams struct {
	ID uuid.UUID `path:"id" doc:"Order ID"`
}

type UpdateOrderStatusInput struct {
	ID   uuid.UUID `path:"id" doc:"Order ID"`
	Body UpdateOrderStatusInputBody
}

type UpdateOrderStatusInputBody struct {
	Status string `json:"status" doc:"Order status" example:"awaiting_pickup"`
}

type UpdateOrderStatusOutput struct {
	Status int
}
