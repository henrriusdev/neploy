package repository

import (
	"context"

	"neploy.dev/pkg/model"
	"neploy.dev/pkg/store"
)

type User interface {
	Create(ctx context.Context, user *model.User) error
	Get(ctx context.Context, id int) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset uint) ([]*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type user[T any] struct {
	Base[T]
}

func NewUser(db store.Queryable) User {
	return &user[model.User]{Base: Base[model.User]{Store: db, Table: "users"}}
}
