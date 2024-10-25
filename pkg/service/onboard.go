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
