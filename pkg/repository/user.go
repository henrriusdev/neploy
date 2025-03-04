package repository

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository/filters"
	"neploy.dev/pkg/store"
)

type User struct {
	Base[model.User]
}

func NewUser(db store.Queryable) *User {
	return &User{Base: Base[model.User]{Store: db, Table: "users"}}
}

func (u *User) Create(ctx context.Context, user model.User) (model.User, error) {
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

func (u *User) Get(ctx context.Context, id string) (model.User, error) {
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

func (u *User) Update(ctx context.Context, user model.User) error {
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

func (u *User) Delete(ctx context.Context, id string) error {
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

func (u *User) List(ctx context.Context, limit, offset uint) ([]model.User, error) {
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

func (u *User) GetByEmail(ctx context.Context, email string) (model.User, error) {
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

func (u *User) CreateInvitation(ctx context.Context, invitation model.Invitation) error {
	query := goqu.Insert("invitations").Rows(invitation)

	q, args, err := query.ToSQL()
	if err != nil {
		return err
	}

	println(q)

	_, err = u.Store.ExecContext(ctx, q, args...)
	return err
}

func (u *User) GetInvitationByToken(ctx context.Context, token string) (model.Invitation, error) {
	query := goqu.From("invitations").Where(goqu.Ex{"token": token})

	q, args, err := query.ToSQL()
	if err != nil {
		return model.Invitation{}, err
	}

	var invitation model.Invitation
	if err := u.Store.QueryRowxContext(ctx, q, args...).StructScan(&invitation); err != nil {
		return model.Invitation{}, err
	}

	return invitation, nil
}

func (u *User) UpdateInvitation(ctx context.Context, invitation model.Invitation) error {
	query := goqu.Update("invitations").
		Set(invitation).
		Where(goqu.Ex{"id": invitation.ID})

	q, args, err := query.ToSQL()
	if err != nil {
		return err
	}

	_, err = u.Store.ExecContext(ctx, q, args...)
	return err
}
