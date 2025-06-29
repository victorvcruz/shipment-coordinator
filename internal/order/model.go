package order

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/victorvcruz/shipment-coordinator/pkg/states"
)

type Status string

const (
	StatusCreated        Status = "created"
	StatusAwaitingPickup Status = "awaiting_pickup"
	StatusPickedUp       Status = "picked_up"
	StatusShipped        Status = "shipped"
	StatusDelivered      Status = "delivered"
	StatusLost           Status = "lost"
)

var StatusValues map[string]Status = map[string]Status{
	"created":         StatusCreated,
	"awaiting_pickup": StatusAwaitingPickup,
	"picked_up":       StatusPickedUp,
	"shipped":         StatusShipped,
	"delivered":       StatusDelivered,
	"lost":            StatusLost,
}

var StatusOrder = map[Status]int{
	StatusCreated:        1,
	StatusAwaitingPickup: 2,
	StatusPickedUp:       3,
	StatusShipped:        4,
	StatusDelivered:      5,
	StatusLost:           6,
}

type Order struct {
	ID            uuid.UUID       `json:"id"`
	Product       string          `json:"product"`
	WeightKg      decimal.Decimal `json:"weight_kg"`
	DestinationUF states.State    `json:"destination_uf"`
	Status        Status          `json:"status"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}
