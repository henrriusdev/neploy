package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository"
)

type Gateway interface {
	Create(ctx context.Context, gateway model.Gateway) error
	Update(ctx context.Context, gateway model.Gateway) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (model.Gateway, error)
	ListByApp(ctx context.Context, appID string) ([]model.Gateway, error)
}

type gateway struct {
	repo repository.Gateway
}

func NewGatewayService(repo repository.Gateway) Gateway {
	return &gateway{
		repo: repo,
	}
}

func (s *gateway) Create(ctx context.Context, gateway model.Gateway) error {
	gateway.ID = uuid.New().String()
	return s.repo.Insert(ctx, gateway)
}

func (s *gateway) Update(ctx context.Context, gateway model.Gateway) error {
	return s.repo.Update(ctx, gateway)
}

func (s *gateway) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *gateway) Get(ctx context.Context, id string) (model.Gateway, error) {
	gateway, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return model.Gateway{}, errors.Wrap(err, "failed to get gateway")
	}
	return gateway, nil
}

func (s *gateway) ListByApp(ctx context.Context, appID string) ([]model.Gateway, error) {
	gateways, err := s.repo.GetByApplicationID(ctx, appID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list gateways")
	}
	return gateways, nil
}
