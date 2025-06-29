package carrier

import (
	"time"

	"github.com/google/uuid"
)

type CreateCarrierInput struct {
	Body CreateCarrierInputBody
}

type CreateCarrierInputBody struct {
	Name     string               `json:"name"     required:"true" doc:"Carrier name"                      example:"Fast Delivery"`
	Policies []CarrierPolicyInput `json:"policies" required:"true" doc:"Region-specific delivery policies"`
}

type CarrierPolicyInput struct {
	Region        string  `json:"region"         required:"true" doc:"Region name"             example:"Nordeste"`
	EstimatedDays int     `json:"estimated_days" required:"true" doc:"Estimated delivery days" example:"5"`
	PricePerKg    float64 `json:"price_per_kg"   required:"true" doc:"Price per kg in BRL"     example:"10.50"`
}

type CarrierResponseOutput struct {
	Status int
	Body   CarrierResponseOutputBody
}

type CarrierResponseOutputBody struct {
	ID        uuid.UUID               `json:"id"         doc:"Carrier ID"               example:"123e4567-e89b-12d3-a456-426614174000"`
	Name      string                  `json:"name"       doc:"Carrier name"             example:"Fast Delivery"`
	Policies  []CarrierPolicyResponse `json:"policies"   doc:"Region-specific policies"`
	CreatedAt time.Time               `json:"created_at" doc:"Carrier creation date"    example:"2023-10-01T12:00:00Z"`
	UpdatedAt time.Time               `json:"updated_at" doc:"Carrier last update date" example:"2023-10-01T12:00:00Z"`
}

type CarrierPolicyResponse struct {
	Region        string    `json:"region"         doc:"Region name"                                 example:"Nordeste"`
	EstimatedDays int       `json:"estimated_days" doc:"Estimated delivery days"                     example:"5"`
	PricePerKg    string    `json:"price_per_kg"   doc:"Price per kg (string for decimal precision)" example:"10.50"`
	CreatedAt     time.Time `json:"created_at"     doc:"Carrier creation date"                       example:"2023-10-01T12:00:00Z"`
	UpdatedAt     time.Time `json:"updated_at"     doc:"Carrier last update date"                    example:"2023-10-01T12:00:00Z"`
}
