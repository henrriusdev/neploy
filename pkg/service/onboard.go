package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/zerolog/log"
	"neploy.dev/pkg/model"
)

type Onboard interface {
	// Onboard the admin user and create the default roles and permissions for the application and also create the users
	Done(context.Context) (bool, int, error)
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

	if _, err := o.metadataService.Get(ctx); err != nil {
		step = 4
		return false, step, err
	}

	step = 5
	return step != 5, step, nil
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
	if _, err := o.roleService.GetByName(ctx, "Administrator"); err != nil && !errors.Is(err, sql.ErrNoRows) {
		role := model.CreateRoleRequest{
			Name:        "Administrator",
			Description: "Administrator of the system",
			Icon:        "User",
			Color:       "#ff0000",
		}
		if err := o.roleService.Create(ctx, role); err != nil {
			log.Err(err).Msg("error creating default role")
		}
	}
	// create the admin user

	req.AdminUser.Roles = []string{"Administrator"}

	if err := o.userService.Create(ctx, req.AdminUser); err != nil {
		log.Err(err).Msg("error users")
		return err
	}

	// create the roles
	for _, role := range req.Roles {
		if err := o.roleService.Create(ctx, role); err != nil {
			log.Err(err).Msg("error roles")
			return err
		}
	}

	// create the metadata
	if err := o.metadataService.Create(ctx, req.Metadata); err != nil {
		log.Err(err).Msg("error meta")
		return err
	}

	return nil
}
