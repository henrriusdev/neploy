package service

import (
	"context"

	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository"
)

type TechStack interface {
	GetAll(context.Context) ([]model.TechStack, error)
	GetByID(context.Context, string) (model.TechStack, error)
	Create(context.Context, model.TechStack) error
	Update(context.Context, model.TechStack) error
	Delete(context.Context, string) error
}

type techStack struct {
	repo repository.TechStack
}

func NewTechStack(repo repository.TechStack) TechStack {
	return &techStack{repo}
}

func (t *techStack) GetAll(ctx context.Context) ([]model.TechStack, error) {
	return t.repo.GetAll(ctx)
}

func (t *techStack) GetByID(ctx context.Context, id string) (model.TechStack, error) {
	return t.repo.GetByID(ctx, id)
}

func (t *techStack) Create(ctx context.Context, techStack model.TechStack) error {
	return t.repo.Insert(ctx, techStack)
}

func (t *techStack) Update(ctx context.Context, techStack model.TechStack) error {
	return t.repo.Update(ctx, techStack)
}

func (t *techStack) Delete(ctx context.Context, id string) error {
	return t.repo.Delete(ctx, id)
}
