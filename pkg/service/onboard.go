package service

import (
	"context"
	"database/sql"
	"errors"

	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
)

type Onboard interface {
	// Onboard the admin user and create the default roles and permissions for the application and also create the users
	Done(context.Context) (bool, error)
	Initiate(context.Context, model.OnboardRequest) error
}

type onboard struct {
	userService     User
	roleService     Role
	metadataService Metadata
}

func NewOnboard(userService User, roleService Role, metadataService Metadata) Onboard {
	return &onboard{userService, roleService, metadataService}
}

func (o *onboard) Done(ctx context.Context) (bool, error) {
	users, err := o.userService.List(ctx, 1, 0)
	if err != nil {
		logger.Error("error getting users: %v", err)
		return false, err
	}

	switch len(users) {
	default:
		return false, nil
	case 1:
		userRoles, err := o.roleService.GetUserRoles(ctx, users[0].ID)
		if err != nil {
			logger.Error("error getting user roles: %v", err)
			return false, err
		}

		if o.hasAdminRole(userRoles) {
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

func (o *onboard) Initiate(ctx context.Context, req model.OnboardRequest) error {
	if _, err := o.roleService.GetByName(ctx, "Administrator"); err != nil && errors.Is(err, sql.ErrNoRows) {
		role := model.CreateRoleRequest{
			Name:        "Administrator",
			Description: "Administrator of the system",
			Icon:        "User",
			Color:       "#ff0000",
		}
		if err := o.roleService.Create(ctx, role); err != nil {
			logger.Error("error creating default role: %v", err)
		}
	}

	if _, err := o.roleService.GetByName(ctx, "Configurator"); err != nil && errors.Is(err, sql.ErrNoRows) {
		role := model.CreateRoleRequest{
			Name:        "Configurator",
			Description: "Gives access to the system configuration for other users without the admin role",
			Icon:        "UserCog",
			Color:       "#000000",
		}
		if err := o.roleService.Create(ctx, role); err != nil {
			logger.Error("error creating default role: %v", err)
		}
	}
	// create the admin user

	req.AdminUser.Roles = []string{"Administrator"}

	if err := o.userService.Create(ctx, req.AdminUser); err != nil {
		logger.Error("error creating admin user: %v", err)
		return err
	}

	// create the roles
	for _, role := range req.Roles {
		if err := o.roleService.Create(ctx, role); err != nil {
			logger.Error("error creating role: %v", err)
			return err
		}
	}

	// create the metadata
	if err := o.metadataService.Create(ctx, req.Metadata); err != nil {
		logger.Error("error creating metadata: %v", err)
		return err
	}

	return nil
}
