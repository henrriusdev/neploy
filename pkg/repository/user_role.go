package repository

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/store"
)

type UserRole interface {
	GetByUserID(context.Context, string) ([]model.UserRoles, error)
	GetByRoleID(context.Context, string) ([]model.UserRoles, error)
	Insert(context.Context, model.UserRoles) (model.UserRoles, error)
}

type userRole[T any] struct {
	Base[T]
}

func NewUserRole(db store.Queryable) UserRole {
	return &userRole[model.UserRoles]{Base[model.UserRoles]{DB: db, Table: "user_roles"}}
}

func (u *userRole[T]) GetByUserID(ctx context.Context, userID string) ([]model.UserRoles, error) {
	q := u.baseQuery("ur").
		Select(
			goqu.I("ur.*"),
			goqu.L(`"u"."id" AS "user.id"`),
			goqu.L(`"u"."name" AS "user.name"`),
			goqu.L(`"u"."email" AS "user.email"`),
			goqu.L(`"r"."id" AS "r.id"`),
			goqu.L(`"r"."name" AS "r.name"`),
			goqu.L(`"r"."description" AS "user.description"`),
		).
		LeftJoin(goqu.T("users").As("u"), goqu.On(goqu.I("u.id"), goqu.I("ur.user_id"))).
		LeftJoin(goqu.T("roles").As("r"), goqu.On(goqu.I("r.id"), goqu.I("ur.role_id"))).
		Where(goqu.I("u.id").Eq(userID))

	query, args, err := q.ToSQL()
	if err != nil {
		return nil, err
	}

	var userRoles []model.UserRoles
	if err := u.Store.SelectContext(ctx, &userRoles, query, args...); err != nil {
		return nil, err
	}

	return userRoles, nil
}

func (u *userRole[T]) GetByRoleID(ctx context.Context, roleID string) ([]model.UserRoles, error) {
	q := u.baseQuery("ur").
		Select(
			goqu.I("ur.*"),
			goqu.L(`"u"."id" AS "user.id"`),
			goqu.L(`"u"."name" AS "user.name"`),
			goqu.L(`"u"."email" AS "user.email"`),
			goqu.L(`"r"."id" AS "r.id"`),
			goqu.L(`"r"."name" AS "r.name"`),
			goqu.L(`"r"."description" AS "user.description"`),
		).
		LeftJoin(goqu.T("users").As("u"), goqu.On(goqu.I("u.id"), goqu.I("ur.user_id"))).
		LeftJoin(goqu.T("roles").As("r"), goqu.On(goqu.I("r.id"), goqu.I("ur.role_id"))).
		Where(goqu.I("r.id").Eq(roleID))

	query, args, err := q.ToSQL()
	if err != nil {
		return nil, err
	}

	var userRoles []model.UserRoles
	if err := u.Store.SelectContext(ctx, &userRoles, query, args...); err != nil {
		return nil, err
	}

	return userRoles, nil
}

func (u *userRole[T]) Insert(ctx context.Context, userRole model.UserRoles) (model.UserRoles, error) {
	q := u.BaseQueryInsert().
		Rows(userRole).
		Returning("*")

	query, args, err := q.ToSQL()
	if err != nil {
		return model.UserRoles{}, err
	}

	if _, err := u.Store.ExecContext(ctx, query, args...); err != nil {
		return model.UserRoles{}, err
	}

	return userRole, nil
}
