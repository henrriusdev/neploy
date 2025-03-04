package service

import (
	"context"

	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository"
)

type Role interface {
	Create(context.Context, model.CreateRoleRequest) error
	GetByID(context.Context, string) (model.Role, error)
	GetByName(context.Context, string) (model.Role, error)
	Get(context.Context) ([]model.RoleWithUsers, error)
	Update(context.Context, string, model.CreateRoleRequest) error
	Delete(context.Context, string) error
	GetUserRoles(context.Context, string) ([]model.UserRoles, error)
}

type role struct {
	roleRepo     *repository.Role
	userRoleRepo *repository.UserRole
}

func NewRole(roleRepo *repository.Role, userRoleRepo *repository.UserRole) Role {
	return &role{roleRepo, userRoleRepo}
}

func (r *role) Create(ctx context.Context, req model.CreateRoleRequest) error {
	role := model.Role{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		Color:       req.Color,
	}
	return r.roleRepo.Insert(ctx, role)
}

func (r *role) GetByID(ctx context.Context, id string) (model.Role, error) {
	return r.roleRepo.GetByID(ctx, id)
}

func (r *role) GetByName(ctx context.Context, name string) (model.Role, error) {
	return r.roleRepo.GetByName(ctx, name)
}

func (r *role) Get(ctx context.Context) ([]model.RoleWithUsers, error) {
	roles, err := r.roleRepo.Get(ctx)
	if err != nil {
		logger.Error("Failed to get roles: %v", err)
		return nil, err
	}

	rolesWithUsers := make([]model.RoleWithUsers, len(roles))
	for i, role := range roles {
		rolesWithUsers[i].Role = role
		userRoles, err := r.userRoleRepo.GetByRoleID(ctx, role.ID)
		if err != nil {
			logger.Error("Failed to get users for role %s: %v", role.Name, err)
			return nil, err
		}

		users := make([]model.User, len(userRoles))
		for i, userRole := range userRoles {
			users[i] = *userRole.User
		}
		rolesWithUsers[i].Users = users
	}
	return rolesWithUsers, nil
}

func (r *role) Update(ctx context.Context, id string, req model.CreateRoleRequest) error {
	role := model.Role{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		Color:       req.Color,
	}
	return r.roleRepo.Update(ctx, id, role)
}

func (r *role) Delete(ctx context.Context, id string) error {
	return r.roleRepo.Delete(ctx, id)
}

func (r *role) GetUserRoles(ctx context.Context, userID string) ([]model.UserRoles, error) {
	return r.userRoleRepo.GetByUserID(ctx, userID)
}
