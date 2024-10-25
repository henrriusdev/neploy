package service

import (
	"context"

	"neploy.dev/pkg/model"
)

type Onboard interface {
	// Onboard the admin user and create the default roles and permissions for the application and also create the users
	Done(context.Context) (bool, error)
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

func (o *onboard) Done(ctx context.Context) (bool, error) {
	users, err := o.userService.List(ctx, 100, 0)
	if err != nil {
		return false, err
	}

	for _, user := range users {
		roles, err := o.roleService.GetUserRoles(ctx, user.ID)
		if err != nil {
			return false, err
		}

		if len(roles) > 0 {
			break
		}

		if o.hasAdminRole(roles) {
			return true, nil
		}
	}

	return false, nil
}

func (o *onboard) hasAdminRole(roles []model.UserRoles) bool {
	for _, role := range roles {
		if role.Role.Name == "Administrator" {
			return true
		}
	}
	return false
}

func (o *onboard) CreateAdminUser(ctx context.Context, req model.CreateUserRequest) error {
	if req.Roles == nil {
		role, err := o.roleService.GetRoleByName(ctx, "Administrator")
		if err != nil {
			return err
		}

		req.Roles = []string{role.ID}
	}

	return o.userService.Create(ctx, req)
}
