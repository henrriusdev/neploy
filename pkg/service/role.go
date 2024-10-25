package service

import (
	"context"

	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository"
)

type Role interface {
	CreateRole(context.Context, model.CreateRoleRequest) error
	GetRoleByID(context.Context, string) (model.Role, error)
	GetRoleByName(context.Context, string) (model.Role, error)
	GetRoles(context.Context) ([]model.Role, error)
	UpdateRole(context.Context, string, model.CreateRoleRequest) error
	DeleteRole(context.Context, string) error
}

type role struct {
	roleRepo repository.Role
}

func NewRole(roleRepo repository.Role) Role {
	return &role{roleRepo}
}

func (r *role) CreateRole(ctx context.Context, req model.CreateRoleRequest) error {
	role := model.Role{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		Color:       req.Color,
	}
	return r.roleRepo.CreateRole(ctx, role)
}

func (r *role) GetRoleByID(ctx context.Context, id string) (model.Role, error) {
	return r.roleRepo.GetRoleByID(ctx, id)
}

func (r *role) GetRoleByName(ctx context.Context, name string) (model.Role, error) {
	return r.roleRepo.GetRoleByName(ctx, name)
}

func (r *role) GetRoles(ctx context.Context) ([]model.Role, error) {
	return r.roleRepo.GetRoles(ctx)
}

func (r *role) UpdateRole(ctx context.Context, id string, req model.CreateRoleRequest) error {
	role := model.Role{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		Color:       req.Color,
	}
	return r.roleRepo.UpdateRole(ctx, id, role)
}

func (r *role) DeleteRole(ctx context.Context, id string) error {
	return r.roleRepo.DeleteRole(ctx, id)
}
