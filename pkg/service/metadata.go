package service

import (
	"context"

	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository"
)

type Metadata interface {
	Create(ctx context.Context, metadata model.MetadataRequest) error
	Update(ctx context.Context, metadata model.MetadataRequest) error
	Get(ctx context.Context) (model.Metadata, error)
	GetTeamName(ctx context.Context) (string, error)
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
		Language: req.Language,
	}
	return m.repo.Create(ctx, metadata)
}

func (m *metadata) Update(ctx context.Context, req model.MetadataRequest) error {
	metadata := model.Metadata{
		TeamName: req.Name,
		LogoURL:  req.LogoURL,
		Language: req.Language,
	}
	return m.repo.Update(ctx, metadata)
}

func (m *metadata) Get(ctx context.Context) (model.Metadata, error) {
	return m.repo.Get(ctx)
}

func (m *metadata) GetTeamName(ctx context.Context) (string, error) {
	return m.repo.GetTeamName(ctx)
}

func (m *metadata) GetTeamLogo(ctx context.Context) (string, error) {
	return m.repo.GetTeamLogo(ctx)
}

func (m *metadata) GetLanguage(ctx context.Context) (string, error) {
	return m.repo.GetLanguage(ctx)
}
