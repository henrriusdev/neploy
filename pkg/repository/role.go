package repository

import (
	"context"

	"neploy.dev/pkg/model"
	"neploy.dev/pkg/store"
)

type Role interface {
	CreateRole(context.Context, model.Role) error
	GetRoleByID(context.Context, string) (model.Role, error)
	GetRoleByName(context.Context, string) (model.Role, error)
	GetRoles(context.Context) ([]model.Role, error)
	UpdateRole(context.Context, string, model.Role) error
	DeleteRole(context.Context, string) error
}

type role[T any] struct {
	Base[T]
}

func NewRole(db store.Queryable) Role {
	return &role[model.Role]{Base[model.Role]{DB: db, Table: "roles"}}
}

func (r *role[T]) CreateRole(ctx context.Context, role model.Role) error {
	return r.Create(ctx, role)
}

func (r *role[T]) GetRoleByID(ctx context.Context, id string) (model.Role, error) {
	var role model.Role
	err := r.GetByID(ctx, id, &role)
	return role, err
}

func (r *role[T]) GetRoleByName(ctx context.Context, name string) (model.Role, error) {
	var role model.Role
	err := r.GetByField(ctx, "name", name, &role)
	return role, err
}

func (r *role[T]) GetRoles(ctx context.Context) ([]model.Role, error) {
	var roles []model.Role
	err := r.Get(ctx, &roles)
	return roles, err
}

func (r *role[T]) UpdateRole(ctx context.Context, id string, role model.Role) error {
	return r.Update(ctx, id, role)
}

func (r *role[T]) DeleteRole(ctx context.Context, id string) error {
	return r.Delete(ctx, id)
}
