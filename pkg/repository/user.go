package repository

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository/filters"
	"neploy.dev/pkg/store"
)

type User interface {
	Create(ctx context.Context, user model.User) (model.User, error)
	Get(ctx context.Context, id string) (model.User, error)
	Update(ctx context.Context, user model.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset uint) ([]model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
}

type user[T any] struct {
	Base[T]
}

func NewUser(db store.Queryable) User {
	return &user[model.User]{Base: Base[model.User]{Store: db, Table: "users"}}
}

func (u *user[T]) Create(ctx context.Context, user model.User) (model.User, error) {
	query := u.BaseQueryInsert().
		Rows(user).
		Returning("*")

	q, args, err := query.ToSQL()
	if err != nil {
		return model.User{}, err
	}

	var newUser model.User
	if err := u.Store.QueryRowxContext(ctx, q, args...).StructScan(&newUser); err != nil {
		return model.User{}, err
	}

	return newUser, nil
}

func (u *user[T]) Get(ctx context.Context, id string) (model.User, error) {
	query := filters.ApplyFilters(u.baseQuery(), filters.IsSelectFilter("id", id))

	q, args, err := query.ToSQL()
	if err != nil {
		return model.User{}, err
	}

	var user model.User
	if err := u.Store.GetContext(ctx, &user, q, args...); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (u *user[T]) Update(ctx context.Context, user model.User) error {
	query := u.BaseQueryUpdate().
		Set(user).
		Where(goqu.Ex{"id": user.ID})

	q, args, err := query.ToSQL()
	if err != nil {
		return err
	}

	if _, err := u.Store.ExecContext(ctx, q, args...); err != nil {
		return err
	}

	return nil
}

func (u *user[T]) Delete(ctx context.Context, id string) error {
	query := u.BaseQueryUpdate().
		Set(goqu.Record{"deleted_at": "CURRENT_TIMESTAMP"}).
		Where(goqu.Ex{"id": id})

	q, args, err := query.ToSQL()
	if err != nil {
		return err
	}

	if _, err := u.Store.ExecContext(ctx, q, args...); err != nil {
		return err
	}

	return nil
}

func (u *user[T]) List(ctx context.Context, limit, offset uint) ([]model.User, error) {
	query := u.baseQuery().
		Limit(limit).
		Offset(offset)

	q, args, err := query.ToSQL()
	if err != nil {
		return nil, err
	}

	var users []model.User
	if err := u.Store.SelectContext(ctx, &users, q, args...); err != nil {
		return nil, err
	}

	return users, nil
}

func (u *user[T]) GetByEmail(ctx context.Context, email string) (model.User, error) {
	query := filters.ApplyFilters(u.baseQuery(), filters.IsSelectFilter("email", email))

	q, args, err := query.ToSQL()
	if err != nil {
		return model.User{}, err
	}

	var user model.User
	if err := u.Store.GetContext(ctx, &user, q, args...); err != nil {
		return model.User{}, err
	}

	return user, nil
}
