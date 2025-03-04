package repository

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository/filters"
	"neploy.dev/pkg/store"
)

type ApplicationUser struct {
	Base[model.ApplicationUser]
}

func NewApplicationUser(db store.Queryable) *ApplicationUser {
	return &ApplicationUser{Base[model.ApplicationUser]{Store: db, Table: "application_users"}}
}

func (a *ApplicationUser) Insert(ctx context.Context, applicationUser model.ApplicationUser) error {
	query := a.BaseQueryInsert().Rows(applicationUser)
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building insert query: %v", err)
		return err
	}

	if _, err := a.Store.ExecContext(ctx, q, args...); err != nil {
		logger.Error("error executing insert query: %v", err)
		return err
	}

	return nil
}

func (a *ApplicationUser) Update(ctx context.Context, applicationUser model.ApplicationUser) error {
	query := filters.ApplyUpdateFilters(a.BaseQueryUpdate().Set(applicationUser), filters.IsUpdateFilter("application_id", applicationUser.ApplicationID), filters.IsUpdateFilter("user_id", applicationUser.UserID))
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building update query: %v", err)
		return err
	}

	if _, err := a.Store.ExecContext(ctx, q, args...); err != nil {
		logger.Error("error executing update query: %v", err)
		return err
	}

	return nil
}

func (a *ApplicationUser) Delete(ctx context.Context, id string) error {
	query := filters.ApplyUpdateFilters(
		a.BaseQueryUpdate().
			Set(goqu.Record{"deleted_at": goqu.L("CURRENT_TIMESTAMP")}),
		filters.IsUpdateFilter("id", id),
	)

	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building delete query: %v", err)
		return err
	}

	if _, err := a.Store.ExecContext(ctx, q, args...); err != nil {
		logger.Error("error executing delete query: %v", err)
		return err
	}

	return nil
}

func (a *ApplicationUser) GetByUserID(ctx context.Context, userID string) ([]model.ApplicationUser, error) {
	query := a.baseQuery().Where(goqu.Ex{"user_id": userID})
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return nil, err
	}

	var applicationUsers []model.ApplicationUser
	if err := a.Store.SelectContext(ctx, &applicationUsers, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return nil, err
	}

	return applicationUsers, nil
}

func (a *ApplicationUser) GetByApplicationID(ctx context.Context, applicationID string) ([]model.ApplicationUser, error) {
	query := a.baseQuery().Where(goqu.Ex{"application_id": applicationID})
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return nil, err
	}

	var applicationUsers []model.ApplicationUser
	if err := a.Store.SelectContext(ctx, &applicationUsers, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return nil, err
	}

	return applicationUsers, nil
}

func (a *ApplicationUser) GetAll(ctx context.Context) ([]model.ApplicationUser, error) {
	query := a.baseQuery().Where(goqu.Ex{"deleted_at": nil})
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return nil, err
	}

	var applicationUsers []model.ApplicationUser
	if err := a.Store.SelectContext(ctx, &applicationUsers, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return nil, err
	}

	return applicationUsers, nil
}
