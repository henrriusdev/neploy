package repository

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/rs/zerolog/log"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository/filters"
	"neploy.dev/pkg/store"
)

type ApplicationUser interface {
	Insert(ctx context.Context, applicationUser model.ApplicationUser) error
	Update(ctx context.Context, applicationUser model.ApplicationUser) error
	Delete(ctx context.Context, id string) error
	GetByUserID(ctx context.Context, userID string) ([]model.ApplicationUser, error)
	GetByApplicationID(ctx context.Context, applicationID string) ([]model.ApplicationUser, error)
	GetAll(ctx context.Context) ([]model.ApplicationUser, error)
}

type applicationUser[T any] struct {
	Base[T]
}

func NewApplicationUser(db store.Queryable) ApplicationUser {
	return &applicationUser[model.ApplicationUser]{Base[model.ApplicationUser]{Store: db, Table: "application_users"}}
}

func (a *applicationUser[T]) Insert(ctx context.Context, applicationUser model.ApplicationUser) error {
	query := a.BaseQueryInsert().Rows(applicationUser)
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building insert query")
		return err
	}

	if _, err := a.Store.ExecContext(ctx, q, args...); err != nil {
		log.Err(err).Msg("error executing insert query")
		return err
	}

	return nil
}

func (a *applicationUser[T]) Update(ctx context.Context, applicationUser model.ApplicationUser) error {
	query := filters.ApplyUpdateFilters(a.BaseQueryUpdate().Set(applicationUser), filters.IsUpdateFilter("application_id", applicationUser.ApplicationID), filters.IsUpdateFilter("user_id", applicationUser.UserID))
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building update query")
		return err
	}

	if _, err := a.Store.ExecContext(ctx, q, args...); err != nil {
		log.Err(err).Msg("error executing update query")
		return err
	}

	return nil
}

func (a *applicationUser[T]) Delete(ctx context.Context, id string) error {
	query := filters.ApplyUpdateFilters(
		a.BaseQueryUpdate().
			Set(goqu.Record{"deleted_at": goqu.L("CURRENT_TIMESTAMP")}),
		filters.IsUpdateFilter("id", id),
	)

	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building delete query")
		return err
	}

	if _, err := a.Store.ExecContext(ctx, q, args...); err != nil {
		log.Err(err).Msg("error executing delete query")
		return err
	}

	return nil
}

func (a *applicationUser[T]) GetByUserID(ctx context.Context, userID string) ([]model.ApplicationUser, error) {
	query := a.baseQuery().Where(goqu.Ex{"user_id": userID})
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building select query")
		return nil, err
	}

	var applicationUsers []model.ApplicationUser
	if err := a.Store.SelectContext(ctx, &applicationUsers, q, args...); err != nil {
		log.Err(err).Msg("error executing select query")
		return nil, err
	}

	return applicationUsers, nil
}

func (a *applicationUser[T]) GetByApplicationID(ctx context.Context, applicationID string) ([]model.ApplicationUser, error) {
	query := a.baseQuery().Where(goqu.Ex{"application_id": applicationID})
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building select query")
		return nil, err
	}

	var applicationUsers []model.ApplicationUser
	if err := a.Store.SelectContext(ctx, &applicationUsers, q, args...); err != nil {
		log.Err(err).Msg("error executing select query")
		return nil, err
	}

	return applicationUsers, nil
}

func (a *applicationUser[T]) GetAll(ctx context.Context) ([]model.ApplicationUser, error) {
	query := a.baseQuery().Where(goqu.Ex{"deleted_at": nil})
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building select query")
		return nil, err
	}

	var applicationUsers []model.ApplicationUser
	if err := a.Store.SelectContext(ctx, &applicationUsers, q, args...); err != nil {
		log.Err(err).Msg("error executing select query")
		return nil, err
	}

	return applicationUsers, nil
}
