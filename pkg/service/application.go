package service

import (
	"context"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository"
)

type Application interface {
	Create(ctx context.Context, app model.Application) error
	Get(ctx context.Context, id string) (model.Application, error)
	GetAll(ctx context.Context) ([]model.Application, error)
	Update(ctx context.Context, app model.Application) error
	GetStat(ctx context.Context, id string) (model.ApplicationStat, error)
	CreateStat(ctx context.Context, stat model.ApplicationStat) error
	UpdateStat(ctx context.Context, stat model.ApplicationStat) error
}

type application struct {
	repo repository.Application
	stat repository.ApplicationStat
}

func NewApplication(repo repository.Application, stat repository.ApplicationStat) Application {
	return &application{repo, stat}
}

func (a *application) Create(ctx context.Context, app model.Application) error {
	return a.repo.Insert(ctx, app)
}

func (a *application) Get(ctx context.Context, id string) (model.Application, error) {
	return a.repo.GetByID(ctx, id)
}

func (a *application) GetAll(ctx context.Context) ([]model.Application, error) {
	return a.repo.GetAll(ctx)
}

func (a *application) Update(ctx context.Context, app model.Application) error {
	return a.repo.Update(ctx, app)
}

func (a *application) GetStat(ctx context.Context, id string) (model.ApplicationStat, error) {
	return a.stat.GetByID(ctx, id)
}

func (a *application) CreateStat(ctx context.Context, stat model.ApplicationStat) error {
	return a.stat.Insert(ctx, stat)
}

func (a *application) UpdateStat(ctx context.Context, stat model.ApplicationStat) error {
	return a.stat.Update(ctx, stat)
}