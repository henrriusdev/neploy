package service

import (
	"context"

	"neploy.dev/pkg/model"
)

type User interface {
	Create(ctx context.Context, user *model.User) error
	Get(ctx context.Context, id int) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset uint) ([]*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type user struct{}
