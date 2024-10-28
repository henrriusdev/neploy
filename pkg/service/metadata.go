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
}

type metadata struct {
	repo repository.Metadata
}

func NewMetadata(repo repository.Metadata) Metadata {
	return &metadata{repo}
}

func (m *metadata) Create(ctx context.Context, req model.MetadataRequest) error {
	metadata := model.Metadata{
		TeamName:       req.Name,
		LogoURL:        req.LogoURL,
		PrimaryColor:   req.PrimaryColor,
		SecondaryColor: req.SecondaryColor,
	}
	return m.repo.Create(ctx, metadata)
}

func (m *metadata) Update(ctx context.Context, req model.MetadataRequest, id string) error {
	metadata := model.Metadata{
		TeamName:       req.Name,
		LogoURL:        req.LogoURL,
		PrimaryColor:   req.PrimaryColor,
		SecondaryColor: req.SecondaryColor,
		BaseEntity:     model.BaseEntity{ID: id},
	}

	return m.repo.Update(ctx, metadata)
}

func (m *metadata) Get(ctx context.Context) (model.Metadata, error) {
	return m.repo.Get(ctx)
}
