package service

import (
	"context"

	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository"
)

type Metadata interface {
	Create(ctx context.Context, metadata model.MetadataRequest) error
	Update(ctx context.Context, metadata model.MetadataRequest, id string) error
	Get(ctx context.Context) (model.Metadata, error)
	GetTeamName(ctx context.Context) (string, error)
	GetPrimaryColor(ctx context.Context) (string, error)
	GetSecondaryColor(ctx context.Context) (string, error)
	GetTeamLogo(ctx context.Context) (string, error)
}

type metadata struct {
	repo repository.Metadata
}

func NewMetadata(repo repository.Metadata) Metadata {
	return &metadata{repo}
}

func (m *metadata) Create(ctx context.Context, req model.MetadataRequest) error {
	metadata := model.Metadata{
		TeamName: req.Name,
		LogoURL:  req.LogoURL,
	}
	return m.repo.Create(ctx, metadata)
}

func (m *metadata) Update(ctx context.Context, req model.MetadataRequest, id string) error {
	metadata := model.Metadata{
		TeamName:   req.Name,
		LogoURL:    req.LogoURL,
		BaseEntity: model.BaseEntity{ID: id},
	}

	return m.repo.Update(ctx, metadata)
}

func (m *metadata) Get(ctx context.Context) (model.Metadata, error) {
	return m.repo.Get(ctx)
}

func (m *metadata) GetPrimaryColor(ctx context.Context) (string, error) {
	return m.repo.GetPrimaryColor(ctx)
}

func (m *metadata) GetSecondaryColor(ctx context.Context) (string, error) {
	return m.repo.GetSecondaryColor(ctx)
}

func (m *metadata) GetTeamName(ctx context.Context) (string, error) {
	return m.repo.GetTeamName(ctx)
}

func (m *metadata) GetTeamLogo(ctx context.Context) (string, error) {
	return m.repo.GetTeamLogo(ctx)
}
