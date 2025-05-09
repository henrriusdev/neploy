package repository

import (
	"context"
	"neploy.dev/pkg/common"

	"github.com/doug-martin/goqu/v9"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository/filters"
	"neploy.dev/pkg/store"
)

type Role struct {
	Base[model.Role]
}

func NewRole(db store.Queryable) *Role {
	return &Role{Base[model.Role]{Store: db, Table: "roles"}}
}

func (r *Role) Insert(ctx context.Context, role model.Role) error {
	q := r.BaseQueryInsert().Rows(role)
	query, args, err := q.ToSQL()
	if err != nil {
		return err
	}

	if _, err := r.Store.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	common.AttachSQLToTrace(ctx, query)
	return nil
}

func (r *Role) GetByID(ctx context.Context, id string) (model.Role, error) {
	var role model.Role
	q := filters.ApplyFilters(r.baseQuery(), filters.IsSelectFilter("id", id))
	query, args, err := q.ToSQL()
	if err != nil {
		return role, err
	}

	if err := r.Store.GetContext(ctx, &role, query, args...); err != nil {
		return role, err
	}

	common.AttachSQLToTrace(ctx, query)
	return role, nil
}

func (r *Role) GetByName(ctx context.Context, name string) (model.Role, error) {
	var role model.Role
	q := filters.ApplyFilters(r.baseQuery(), filters.IsSelectFilter("name", name))
	query, args, err := q.ToSQL()
	if err != nil {
		return role, err
	}

	if err := r.Store.GetContext(ctx, &role, query, args...); err != nil {
		return role, err
	}

	common.AttachSQLToTrace(ctx, query)
	return role, nil
}

func (r *Role) Get(ctx context.Context) ([]model.Role, error) {
	var roles []model.Role
	q := r.baseQuery()
	query, args, err := q.ToSQL()
	if err != nil {
		return roles, err
	}

	if err := r.Store.SelectContext(ctx, &roles, query, args...); err != nil {
		return roles, err
	}

	common.AttachSQLToTrace(ctx, query)
	return roles, nil
}

func (r *Role) Update(ctx context.Context, id string, role model.Role) error {
	q := r.BaseQueryUpdate().Set(role).Where(goqu.Ex{"id": id})
	query, args, err := q.ToSQL()
	if err != nil {
		return err
	}

	if _, err := r.Store.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	common.AttachSQLToTrace(ctx, query)
	return nil
}

func (r *Role) Delete(ctx context.Context, id string) error {
	q := r.BaseQueryUpdate().Where(goqu.Ex{"id": id}).Set(goqu.Record{"deleted_at": goqu.L("CURRENT_TIMESTAMP")})
	query, args, err := q.ToSQL()
	if err != nil {
		return err
	}

	if _, err := r.Store.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	common.AttachSQLToTrace(ctx, query)
	return nil
}
