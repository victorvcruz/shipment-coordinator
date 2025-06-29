package carrier

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/victorvcruz/shipment-coordinator/pkg/states"
)

type Carrier struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Policies  []Policy  `json:"policies"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Policy struct {
	ID            uuid.UUID       `json:"id"`
	CarrierID     uuid.UUID       `json:"carrier_id"`
	Region        states.Region   `json:"region"`
	EstimatedDays int             `json:"estimated_days"`
	PricePerKg    decimal.Decimal `json:"price_per_kg"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}
