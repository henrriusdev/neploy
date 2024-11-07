package repository

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/rs/zerolog/log"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository/filters"
	"neploy.dev/pkg/store"
)

type Gateway interface {
	Insert(ctx context.Context, gateway model.Gateway) error
	Update(ctx context.Context, gateway model.Gateway) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (model.Gateway, error)
	GetAll(ctx context.Context) ([]model.Gateway, error)
	GetByHttpMethod(ctx context.Context, httpMethod string) ([]model.Gateway, error)
	GetByEndpoint(ctx context.Context, endpoint string) ([]model.Gateway, error)
	GetByLogLevel(ctx context.Context, logLevel string) ([]model.Gateway, error)
	GetByStage(ctx context.Context, stage string) ([]model.Gateway, error)
	GetByName(ctx context.Context, name string) ([]model.Gateway, error)
	GetByApplicationID(ctx context.Context, applicationID string) ([]model.Gateway, error)
}

type gateway[T any] struct {
	Base[T]
}

func NewGateway(db store.Queryable) Gateway {
	return &gateway[model.Gateway]{Base[model.Gateway]{Store: db, Table: "gateways"}}
}

func (g *gateway[T]) Insert(ctx context.Context, gateway model.Gateway) error {
	query := g.BaseQueryInsert().Rows(gateway)
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building insert query")
		return err
	}

	if _, err := g.Store.ExecContext(ctx, q, args...); err != nil {
		log.Err(err).Msg("error executing insert query")
		return err
	}

	return nil
}

func (g *gateway[T]) Update(ctx context.Context, gateway model.Gateway) error {
	query := filters.ApplyUpdateFilters(g.BaseQueryUpdate().Set(gateway), filters.IsUpdateFilter("id", gateway.ID))
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building update query")
		return err
	}

	if _, err := g.Store.ExecContext(ctx, q, args...); err != nil {
		log.Err(err).Msg("error executing update query")
		return err
	}

	return nil
}

func (g *gateway[T]) Delete(ctx context.Context, id string) error {
	query := filters.ApplyUpdateFilters(
		g.BaseQueryUpdate().
			Set(goqu.Record{"deleted_at": goqu.L("CURRENT_TIMESTAMP")}),
		filters.IsUpdateFilter("id", id),
	)

	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building delete query")
		return err
	}

	if _, err := g.Store.ExecContext(ctx, q, args...); err != nil {
		log.Err(err).Msg("error executing delete query")
		return err
	}

	return nil
}

func (g *gateway[T]) GetByID(ctx context.Context, id string) (model.Gateway, error) {
	query := g.baseQuery().Where(goqu.Ex{"id": id})
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building select query")
		return model.Gateway{}, err
	}

	var gateway model.Gateway
	if err := g.Store.GetContext(ctx, &gateway, q, args...); err != nil {
		log.Err(err).Msg("error executing select query")
		return model.Gateway{}, err
	}

	return gateway, nil
}

func (g *gateway[T]) GetAll(ctx context.Context) ([]model.Gateway, error) {
	query := g.baseQuery()
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building select query")
		return nil, err
	}

	var gateways []model.Gateway
	if err := g.Store.SelectContext(ctx, &gateways, q, args...); err != nil {
		log.Err(err).Msg("error executing select query")
		return nil, err
	}

	return gateways, nil
}

func (g *gateway[T]) GetByHttpMethod(ctx context.Context, httpMethod string) ([]model.Gateway, error) {
	query := g.baseQuery().Where(goqu.Ex{"http_method": httpMethod})
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building select query")
		return nil, err
	}

	var gateways []model.Gateway
	if err := g.Store.SelectContext(ctx, &gateways, q, args...); err != nil {
		log.Err(err).Msg("error executing select query")
		return nil, err
	}

	return gateways, nil
}

func (g *gateway[T]) GetByEndpoint(ctx context.Context, endpoint string) ([]model.Gateway, error) {
	query := g.baseQuery().Where(goqu.Ex{"endpoint": endpoint})
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building select query")
		return nil, err
	}

	var gateways []model.Gateway
	if err := g.Store.SelectContext(ctx, &gateways, q, args...); err != nil {
		log.Err(err).Msg("error executing select query")
		return nil, err
	}

	return gateways, nil
}

func (g *gateway[T]) GetByLogLevel(ctx context.Context, logLevel string) ([]model.Gateway, error) {
	query := g.baseQuery().Where(goqu.Ex{"log_level": logLevel})
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building select query")
		return nil, err
	}

	var gateways []model.Gateway
	if err := g.Store.SelectContext(ctx, &gateways, q, args...); err != nil {
		log.Err(err).Msg("error executing select query")
		return nil, err
	}

	return gateways, nil
}

func (g *gateway[T]) GetByStage(ctx context.Context, stage string) ([]model.Gateway, error) {
	query := g.baseQuery().Where(goqu.Ex{"stage": stage})
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building select query")
		return nil, err
	}

	var gateways []model.Gateway
	if err := g.Store.SelectContext(ctx, &gateways, q, args...); err != nil {
		log.Err(err).Msg("error executing select query")
		return nil, err
	}

	return gateways, nil
}

func (g *gateway[T]) GetByName(ctx context.Context, name string) ([]model.Gateway, error) {
	query := g.baseQuery().Where(goqu.Ex{"name": name})
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building select query")
		return nil, err
	}

	var gateways []model.Gateway
	if err := g.Store.SelectContext(ctx, &gateways, q, args...); err != nil {
		log.Err(err).Msg("error executing select query")
		return nil, err
	}

	return gateways, nil
}

func (g *gateway[T]) GetByApplicationID(ctx context.Context, applicationID string) ([]model.Gateway, error) {
	query := g.baseQuery().Where(goqu.Ex{"application_id": applicationID})
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building select query")
		return nil, err
	}

	var gateways []model.Gateway
	if err := g.Store.SelectContext(ctx, &gateways, q, args...); err != nil {
		log.Err(err).Msg("error executing select query")
		return nil, err
	}

	return gateways, nil
}
