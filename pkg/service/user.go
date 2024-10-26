package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository"
)

type User interface {
	Create(ctx context.Context, user model.CreateUserRequest) error
	Get(ctx context.Context, id string) (model.User, error)
	Update(ctx context.Context, user model.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset uint) ([]model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
}

type user struct {
	repo         repository.User
	userRoleRepo repository.UserRole
}

func NewUser(repo repository.User, userRoleRepo repository.UserRole) User {
	return &user{repo, userRoleRepo}
}

func (u *user) Create(ctx context.Context, req model.CreateUserRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := model.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Address:   req.Address,
		Phone:     req.Phone,
	}

	user, err = u.repo.Create(ctx, user)
	if err != nil {
		return err
	}

	for _, roleID := range req.Roles {
		userRole := model.UserRoles{
			UserID: user.ID,
			RoleID: roleID,
		}
		if _, err := u.userRoleRepo.Insert(ctx, userRole); err != nil {
			return err
		}
	}

	return nil
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
