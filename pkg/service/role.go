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
