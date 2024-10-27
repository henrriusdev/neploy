package service

import (
	"context"

	"neploy.dev/pkg/model"
)

type Onboard interface {
	// Onboard the admin user and create the default roles and permissions for the application and also create the users
	Done(context.Context) (bool, int, error)
	Initiate(context.Context, model.OnboardRequest) error
}

type onboard struct {
	userService User
	roleService Role
}

func NewOnboard(userService User, roleService Role) Onboard {
	return &onboard{userService, roleService}
}

func (o *onboard) Done(ctx context.Context) (bool, int, error) {
	users, err := o.userService.List(ctx, 100, 0)
	if err != nil {
		return false, 0, err
	}

	step := 0
	switch len(users) {
	case 0:
		return false, 0, nil
	case 1:
		userRoles, err := o.roleService.GetUserRoles(ctx, users[0].ID)
		if err != nil {
			return false, 0, err
		}

		if o.hasAdminRole(userRoles) {
			step = 1
		}

		roles, err := o.roleService.Get(ctx)
		if err != nil {
			return false, step, err
		}

		if len(roles) == 0 {
			return false, step, nil
		}

		step = 2
	default:
		step = 3
	}

	// todo: check if the metadata is created

	return step != 0, step, nil
}

func (o *onboard) hasAdminRole(roles []model.UserRoles) bool {
	for _, role := range roles {
		if role.Role.Name == "Administrator" {
			return true
		}
	}
	return false
}

func (o *onboard) Initiate(ctx context.Context, req model.OnboardRequest) error {
	// create the admin user
	if err := o.userService.Create(ctx, req.AdminUser); err != nil {
		return err
	}

	// create the roles
	for _, role := range req.Roles {
		if err := o.roleService.Create(ctx, role); err != nil {
			return err
		}
	}

	// create the users
	for _, user := range req.Users {
		if err := o.userService.Create(ctx, user); err != nil {
			return err
		}
	}

	// todo: create the metadata

	return nil
}
