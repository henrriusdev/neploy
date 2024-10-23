package service

import (
	"context"

	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository"
)

type User interface {
	Create(ctx context.Context, user model.User) error
	Get(ctx context.Context, id string) (model.User, error)
	Update(ctx context.Context, user model.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset uint) ([]model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
}

type user struct {
	repo repository.User
}

func NewUser(repo repository.User) User {
	return &user{repo}
}

func (u *user) Create(ctx context.Context, user model.User) error {
	return u.repo.Create(ctx, user)
}

func (u *user) Get(ctx context.Context, id string) (model.User, error) {
	return u.repo.Get(ctx, id)
}

func (u *user) Update(ctx context.Context, user model.User) error {
	return u.repo.Update(ctx, user)
}

func (u *user) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}

func (u *user) List(ctx context.Context, limit, offset uint) ([]model.User, error) {
	return u.repo.List(ctx, limit, offset)
}

func (u *user) GetByEmail(ctx context.Context, email string) (model.User, error) {
	return u.repo.GetByEmail(ctx, email)
}
