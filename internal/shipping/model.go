package shipping

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Contract struct {
	ID            uuid.UUID       `json:"id"`
	OrderID       uuid.UUID       `json:"order_id"`
	CarrierID     uuid.UUID       `json:"carrier"`
	Price         decimal.Decimal `json:"price"`
	EstimatedDays int             `json:"estimated_days"`
	ContractedAt  time.Time       `json:"contracted_at"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type Quote struct {
	CarrierID     uuid.UUID       `json:"carrier_id"`
	CarrierName   string          `json:"carrier_name"`
	Price         decimal.Decimal `json:"price"`
	EstimatedDays int             `json:"estimated_days"`
}
