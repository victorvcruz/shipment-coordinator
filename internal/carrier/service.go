package carrier

import (
	"context"

	log "go.uber.org/zap"
)

//go:generate moq -pkg mocks -out mocks/service.go . Service
type Service interface {
	Create(ctx context.Context, carrier *Carrier) (*Carrier, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, carrier *Carrier) (*Carrier, error) {
	id, err := s.repo.Create(ctx, carrier)
	if err != nil {
		log.L().
			Error("failed to create carrier", log.String("carrier_name", carrier.Name), log.Error(err))
		return nil, err
	}

	carrier.ID = id
	return carrier, nil
}
