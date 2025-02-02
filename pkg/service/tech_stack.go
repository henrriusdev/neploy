package service

import (
	"context"

	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository"
)

type TechStack interface {
	GetAll(context.Context) ([]model.TechStackWithApplications, error)
	GetByID(context.Context, string) (model.TechStack, error)
	Create(context.Context, model.CreateTechStackRequest) error
	Update(context.Context, string, model.CreateTechStackRequest) error
	Delete(context.Context, string) error
}

type techStack struct {
	repo    repository.TechStack
	appRepo repository.Application
}

func NewTechStack(repo repository.TechStack, appRepo repository.Application) TechStack {
	return &techStack{repo, appRepo}
}

func (t *techStack) GetAll(ctx context.Context) ([]model.TechStackWithApplications, error) {
	stacks, err := t.repo.GetAll(ctx)
	if err != nil {
		logger.Error("error getting tech stacks: %v", err)
		return nil, err
	}

	techStacks := make([]model.TechStackWithApplications, len(stacks))
	for i, stack := range stacks {
		techStacks[i].TechStack = stack
		apps, err := t.appRepo.GetByTechStack(ctx, stack.ID)
		if err != nil {
			logger.Error("error getting applications by tech stack: %v", err)
			return nil, err
		}

		techStacks[i].Applications = apps
	}

	return techStacks, nil
}

func (t *techStack) GetByID(ctx context.Context, id string) (model.TechStack, error) {
	return t.repo.GetByID(ctx, id)
}

func (t *techStack) Create(ctx context.Context, techStack model.CreateTechStackRequest) error {
	tech := model.TechStack{
		Name:        techStack.Name,
		Description: techStack.Description,
	}
	return t.repo.Insert(ctx, tech)
}

func (t *techStack) Update(ctx context.Context, id string, techStack model.CreateTechStackRequest) error {
	tech := model.TechStack{
		Name:        techStack.Name,
		Description: techStack.Description,
	}
	return t.repo.Update(ctx, id, tech)
}

func (t *techStack) Delete(ctx context.Context, id string) error {
	return t.repo.Delete(ctx, id)
}
