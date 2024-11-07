package repository

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/rs/zerolog/log"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository/filters"
	"neploy.dev/pkg/store"
)

type Environment interface {
	Insert(ctx context.Context, environment model.Environment) error
	Update(ctx context.Context, environment model.Environment) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (model.Environment, error)
	GetAll(ctx context.Context) ([]model.Environment, error)
	GetByName(ctx context.Context, name string) (model.Environment, error)
}

type environment[T any] struct {
	Base[T]
}

func NewEnvironment(db store.Queryable) Environment {
	return &environment[model.Environment]{Base[model.Environment]{Store: db, Table: "environments"}}
}

func (e *environment[T]) Insert(ctx context.Context, environment model.Environment) error {
	query := e.BaseQueryInsert().Rows(environment)
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building insert query")
		return err
	}

	if _, err := e.Store.ExecContext(ctx, q, args...); err != nil {
		log.Err(err).Msg("error executing insert query")
		return err
	}

	return nil
}

func (e *environment[T]) Update(ctx context.Context, environment model.Environment) error {
	query := filters.ApplyUpdateFilters(e.BaseQueryUpdate().Set(environment), filters.IsUpdateFilter("id", environment.ID))
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building update query")
		return err
	}

	if _, err := e.Store.ExecContext(ctx, q, args...); err != nil {
		log.Err(err).Msg("error executing update query")
		return err
	}

	return nil
}

func (e *environment[T]) Delete(ctx context.Context, id string) error {
	query := filters.ApplyUpdateFilters(
		e.BaseQueryUpdate().
			Set(goqu.Record{"deleted_at": goqu.L("CURRENT_TIMESTAMP")}),
		filters.IsUpdateFilter("id", id),
	)

	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building delete query")
		return err
	}

	if _, err := e.Store.ExecContext(ctx, q, args...); err != nil {
		log.Err(err).Msg("error executing delete query")
		return err
	}

	return nil
}

func (e *environment[T]) GetByID(ctx context.Context, id string) (model.Environment, error) {
	query := e.baseQuery().Where(goqu.Ex{"id": id})
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building get by id query")
		return model.Environment{}, err
	}

	var environment model.Environment
	if err := e.Store.GetContext(ctx, &environment, q, args...); err != nil {
		log.Err(err).Msg("error executing get by id query")
		return model.Environment{}, err
	}

	return environment, nil
}

func (e *environment[T]) GetAll(ctx context.Context) ([]model.Environment, error) {
	query := e.baseQuery()
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building select query")
		return nil, err
	}

	var environments []model.Environment
	if err := e.Store.SelectContext(ctx, &environments, q, args...); err != nil {
		log.Err(err).Msg("error executing select query")
		return nil, err
	}

	return environments, nil
}

func (e *environment[T]) GetByName(ctx context.Context, name string) (model.Environment, error) {
	query := e.baseQuery().Where(goqu.Ex{"name": name})
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building get by name query")
		return model.Environment{}, err
	}

	var environment model.Environment
	if err := e.Store.GetContext(ctx, &environment, q, args...); err != nil {
		log.Err(err).Msg("error executing get by name query")
		return model.Environment{}, err
	}

	return environment, nil
}
