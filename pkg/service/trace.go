package service

import (
	"context"

	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository"
)

type Trace interface {
	GetAll(context.Context) ([]model.Trace, error)
	GetByID(context.Context, string) (model.Trace, error)
	Create(context.Context, model.Trace) error
	Update(context.Context, model.Trace) error
	Delete(context.Context, string) error
}

type trace struct {
	repo *repository.Trace
}

func NewTrace(repo *repository.Trace) Trace {
	return &trace{repo}
}

func (t *trace) GetAll(ctx context.Context) ([]model.Trace, error) {
	return t.repo.GetAll(ctx)
}

func (t *trace) GetByID(ctx context.Context, id string) (model.Trace, error) {
	return t.repo.GetByID(ctx, id)
}

func (t *trace) Create(ctx context.Context, trace model.Trace) error {
	return t.repo.Insert(ctx, trace)
}

func (t *trace) Update(ctx context.Context, trace model.Trace) error {
	return t.repo.Update(ctx, trace)
}

func (t *trace) Delete(ctx context.Context, id string) error {
	return t.repo.Delete(ctx, id)
}
