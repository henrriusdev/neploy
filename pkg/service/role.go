package service

import (
	"context"

	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository"
)

type Role interface {
	Create(context.Context, model.CreateRoleRequest) error
	GetByID(context.Context, string) (model.Role, error)
	GetByName(context.Context, string) (model.Role, error)
	Get(context.Context) ([]model.Role, error)
	Update(context.Context, string, model.CreateRoleRequest) error
	Delete(context.Context, string) error
	GetUserRoles(context.Context, string) ([]model.UserRoles, error)
}

type role struct {
	roleRepo     repository.Role
	userRoleRepo repository.UserRole
}

func NewRole(roleRepo repository.Role, userRoleRepo repository.UserRole) Role {
	return &role{roleRepo, userRoleRepo}
}

func (r *role) Create(ctx context.Context, req model.CreateRoleRequest) error {
	role := model.Role{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		Color:       req.Color,
	}
	return r.roleRepo.CreateRole(ctx, role)
}

func (r *role) GetByID(ctx context.Context, id string) (model.Role, error) {
	return r.roleRepo.GetRoleByID(ctx, id)
}

func (r *role) GetByName(ctx context.Context, name string) (model.Role, error) {
	return r.roleRepo.GetRoleByName(ctx, name)
}

func (r *role) Get(ctx context.Context) ([]model.Role, error) {
	return r.roleRepo.GetRoles(ctx)
}

func (r *role) Update(ctx context.Context, id string, req model.CreateRoleRequest) error {
	role := model.Role{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		Color:       req.Color,
	}
	return r.roleRepo.UpdateRole(ctx, id, role)
}

func (r *role) Delete(ctx context.Context, id string) error {
	return r.roleRepo.DeleteRole(ctx, id)
}

func (r *role) GetUserRoles(ctx context.Context, userID string) ([]model.UserRoles, error) {
	return r.userRoleRepo.GetByUserID(ctx, userID)
}
