package shipping

import (
	"time"

	"github.com/google/uuid"
)

type GetQuotesInput struct {
	OrderID string `path:"order_id" doc:"Order ID"`
}

type QuotesOutput struct {
	Status int
	Body   []QuotesOutputBody
}

type QuotesOutputBody struct {
	CarrierID     string `json:"carrier_id"     doc:"Carrier ID"              example:"123e4567-e89b-12d3-a456-426614174000"`
	CarrierName   string `json:"carrier_name"   doc:"Carrier name"            example:"Fast Delivery"`
	Price         string `json:"price"          doc:"Price in BRL"            example:"10.50"`
	EstimatedDays int    `json:"estimated_days" doc:"Estimated delivery days" example:"5"`
}

type ContractCarrierInput struct {
	Body ContractCarrierInputBody
}

type ContractCarrierInputBody struct {
	OrderID   uuid.UUID `json:"order_id"   required:"true" doc:"Order ID"   example:"111e4567-e89b-12d3-a456-426614174000"`
	CarrierID uuid.UUID `json:"carrier_id" required:"true" doc:"Carrier ID" example:"222e4567-e89b-12d3-a456-426614174000"`
}

type ContractCarrierOutput struct {
	Status int                       `json:"status"`
	Body   ContractCarrierOutputBody `json:"body"`
}

type ContractCarrierOutputBody struct {
	ID            string    `json:"id"             doc:"Contract ID"                 example:"123e4567-e89b-12d3-a456-426614174000"`
	OrderID       string    `json:"order_id"       doc:"Order ID"                    example:"111e4567-e89b-12d3-a456-426614174000"`
	CarrierID     string    `json:"carrier_id"     doc:"Carrier ID"                  example:"222e4567-e89b-12d3-a456-426614174000"`
	Price         string    `json:"price"          doc:"Total price"                 example:"25.50"`
	EstimatedDays int       `json:"estimated_days" doc:"Delivery estimation in days" example:"4"`
	ContractedAt  time.Time `json:"contracted_at"  doc:"Contract date"               example:"2025-06-28T15:04:05Z"`
	CreatedAt     time.Time `json:"created_at"     doc:"Creation timestamp"          example:"2025-06-28T15:04:05Z"`
	UpdatedAt     time.Time `json:"updated_at"     doc:"Last update timestamp"       example:"2025-06-28T15:04:05Z"`
}
