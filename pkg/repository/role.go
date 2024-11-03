package repository

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository/filters"
	"neploy.dev/pkg/store"
)

type Role interface {
	Insert(context.Context, model.Role) error
	GetByID(context.Context, string) (model.Role, error)
	GetByName(context.Context, string) (model.Role, error)
	Get(context.Context) ([]model.Role, error)
	Update(context.Context, string, model.Role) error
	Delete(context.Context, string) error
}

type role[T any] struct {
	Base[T]
}

func NewRole(db store.Queryable) Role {
	return &role[model.Role]{Base[model.Role]{Store: db, Table: "roles"}}
}

func (r *role[T]) Insert(ctx context.Context, role model.Role) error {
	q := r.BaseQueryInsert().Rows(role)
	query, args, err := q.ToSQL()
	if err != nil {
		return err
	}

	if _, err := r.Store.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *role[T]) GetByID(ctx context.Context, id string) (model.Role, error) {
	var role model.Role
	q := filters.ApplyFilters(r.baseQuery(), filters.IsSelectFilter("id", id))
	query, args, err := q.ToSQL()
	if err != nil {
		return role, err
	}

	if err := r.Store.GetContext(ctx, &role, query, args...); err != nil {
		return role, err
	}

	return role, nil
}

func (r *role[T]) GetByName(ctx context.Context, name string) (model.Role, error) {
	var role model.Role
	q := filters.ApplyFilters(r.baseQuery(), filters.IsSelectFilter("name", name))
	query, args, err := q.ToSQL()
	if err != nil {
		return role, err
	}

	if err := r.Store.GetContext(ctx, &role, query, args...); err != nil {
		return role, err
	}

	return role, nil
}

func (r *role[T]) Get(ctx context.Context) ([]model.Role, error) {
	var roles []model.Role
	q := r.baseQuery()
	query, args, err := q.ToSQL()
	if err != nil {
		return roles, err
	}

	if err := r.Store.SelectContext(ctx, &roles, query, args...); err != nil {
		return roles, err
	}

	return roles, nil
}

func (r *role[T]) Update(ctx context.Context, id string, role model.Role) error {
	q := r.BaseQueryUpdate().Set(role).Where(goqu.Ex{"id": id})
	query, args, err := q.ToSQL()
	if err != nil {
		return err
	}

	if _, err := r.Store.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *role[T]) Delete(ctx context.Context, id string) error {
	role, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	role.DeletedAt = &model.Date{Time: time.Now()}
	return r.Update(ctx, id, role)
}
