package repository

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/rs/zerolog/log"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository/filters"
	"neploy.dev/pkg/store"
)

type ApplicationEnvironment interface {
	Insert(ctx context.Context, applicationEnvironment model.ApplicationEnvironment) error
	Update(ctx context.Context, applicationEnvironment model.ApplicationEnvironment) error
	Delete(ctx context.Context, id string) error
	GetByApplicationID(ctx context.Context, applicationID string) ([]model.ApplicationEnvironment, error)
	GetByEnvironmentID(ctx context.Context, environmentID string) ([]model.ApplicationEnvironment, error)
	GetByApplicationIDAndEnvironmentID(ctx context.Context, applicationID, environmentID string) (model.ApplicationEnvironment, error)
	GetAll(ctx context.Context) ([]model.ApplicationEnvironment, error)
}

type applicationEnvironment[T any] struct {
	Base[T]
}

func NewApplicationEnvironment(db store.Queryable) ApplicationEnvironment {
	return &applicationEnvironment[model.ApplicationEnvironment]{Base[model.ApplicationEnvironment]{Store: db, Table: "application_environments"}}
}

func (a *applicationEnvironment[T]) Insert(ctx context.Context, applicationEnvironment model.ApplicationEnvironment) error {
	query := a.BaseQueryInsert().Rows(applicationEnvironment)
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

func (a *applicationEnvironment[T]) Update(ctx context.Context, applicationEnvironment model.ApplicationEnvironment) error {
	query := filters.ApplyUpdateFilters(
		a.BaseQueryUpdate().Set(applicationEnvironment),
		filters.IsUpdateFilter("application_id", applicationEnvironment.ApplicationID),
		filters.IsUpdateFilter("environment_id", applicationEnvironment.EnvironmentID),
	)
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

func (a *applicationEnvironment[T]) Delete(ctx context.Context, id string) error {
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

func (a *applicationEnvironment[T]) GetByApplicationID(ctx context.Context, applicationID string) ([]model.ApplicationEnvironment, error) {
	query := a.baseQuery().Where(goqu.Ex{"application_id": applicationID, "deleted_at": nil})
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building get by application id query")
		return nil, err
	}

	var applicationEnvironments []model.ApplicationEnvironment
	if err := a.Store.SelectContext(ctx, &applicationEnvironments, q, args...); err != nil {
		log.Err(err).Msg("error executing get by application id query")
		return nil, err
	}

	return applicationEnvironments, nil
}

func (a *applicationEnvironment[T]) GetByEnvironmentID(ctx context.Context, environmentID string) ([]model.ApplicationEnvironment, error) {
	query := a.baseQuery().Where(goqu.Ex{"environment_id": environmentID, "deleted_at": nil})
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building get by environment id query")
		return nil, err
	}

	var applicationEnvironments []model.ApplicationEnvironment
	if err := a.Store.SelectContext(ctx, &applicationEnvironments, q, args...); err != nil {
		log.Err(err).Msg("error executing get by environment id query")
		return nil, err
	}

	return applicationEnvironments, nil
}

func (a *applicationEnvironment[T]) GetByApplicationIDAndEnvironmentID(ctx context.Context, applicationID, environmentID string) (model.ApplicationEnvironment, error) {
	query := a.baseQuery().Where(goqu.Ex{"application_id": applicationID, "environment_id": environmentID, "deleted_at": nil})
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building get by application id and environment id query")
		return model.ApplicationEnvironment{}, err
	}

	var applicationEnvironment model.ApplicationEnvironment
	if err := a.Store.GetContext(ctx, &applicationEnvironment, q, args...); err != nil {
		log.Err(err).Msg("error executing get by application id and environment id query")
		return model.ApplicationEnvironment{}, err
	}

	return applicationEnvironment, nil
}

func (a *applicationEnvironment[T]) GetAll(ctx context.Context) ([]model.ApplicationEnvironment, error) {
	query := a.baseQuery().Where(goqu.Ex{"deleted_at": nil})
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building select query")
		return nil, err
	}

	var applicationEnvironments []model.ApplicationEnvironment
	if err := a.Store.SelectContext(ctx, &applicationEnvironments, q, args...); err != nil {
		log.Err(err).Msg("error executing select query")
		return nil, err
	}

	return applicationEnvironments, nil
}
