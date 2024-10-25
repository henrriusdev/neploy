package service

import (
	"context"

	"neploy.dev/pkg/model"
)

type Onboard interface {
	// Onboard the admin user and create the default roles and permissions for the application and also create the users
	Done(context.Context) error
	CreateAdminUser(context.Context, model.CreateUserRequest) error
	CreateRole(context.Context, model.CreateRoleRequest) error
}

type onboard struct {
	userService User
	roleService Role
}

func NewOnboard(userService User, roleService Role) Onboard {
	return &onboard{userService, roleService}
}

func (o *onboard) Done(ctx context.Context) error {
	users, err := o.userService.List(ctx, 100, 0)
	if err != nil {
		return err
	}

	for _, user := range users {
	}
}
