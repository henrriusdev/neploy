package repository

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository/filters"
	"neploy.dev/pkg/store"
)

type Gateway struct {
	Base[model.Gateway]
}

func NewGateway(db store.Queryable) *Gateway {
	return &Gateway{Base[model.Gateway]{Store: db, Table: "gateways"}}
}

func (g *Gateway) Insert(ctx context.Context, gateway model.Gateway) error {
	query := g.BaseQueryInsert().Rows(gateway)
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building insert query: %v", err)
		return err
	}

	if _, err := g.Store.ExecContext(ctx, q, args...); err != nil {
		logger.Error("error executing insert query: %v", err)
		return err
	}

	return nil
}

func (g *Gateway) Update(ctx context.Context, gateway model.Gateway) error {
	query := filters.ApplyUpdateFilters(g.BaseQueryUpdate().Set(gateway), filters.IsUpdateFilter("id", gateway.ID))
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building update query: %v", err)
		return err
	}

	if _, err := g.Store.ExecContext(ctx, q, args...); err != nil {
		logger.Error("error executing update query: %v", err)
		return err
	}

	return nil
}

func (g *Gateway) Delete(ctx context.Context, id string) error {
	query := filters.ApplyUpdateFilters(
		g.BaseQueryUpdate().
			Set(goqu.Record{"deleted_at": goqu.L("CURRENT_TIMESTAMP")}),
		filters.IsUpdateFilter("id", id),
	)

	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building delete query: %v", err)
		return err
	}

	if _, err := g.Store.ExecContext(ctx, q, args...); err != nil {
		logger.Error("error executing delete query: %v", err)
		return err
	}

	return nil
}

func (g *Gateway) GetByID(ctx context.Context, id string) (model.Gateway, error) {
	query := g.baseQuery().Where(goqu.Ex{"id": id})
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return model.Gateway{}, err
	}

	var gateway model.Gateway
	if err := g.Store.GetContext(ctx, &gateway, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return model.Gateway{}, err
	}

	return gateway, nil
}

func (g *Gateway) GetAll(ctx context.Context) ([]model.Gateway, error) {
	query := g.baseQuery()
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return nil, err
	}

	var gateways []model.Gateway
	if err := g.Store.SelectContext(ctx, &gateways, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return nil, err
	}

	return gateways, nil
}

func (g *Gateway) GetByHttpMethod(ctx context.Context, httpMethod string) ([]model.Gateway, error) {
	query := g.baseQuery().Where(goqu.Ex{"http_method": httpMethod})
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return nil, err
	}

	var gateways []model.Gateway
	if err := g.Store.SelectContext(ctx, &gateways, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return nil, err
	}

	return gateways, nil
}

func (g *Gateway) GetByEndpoint(ctx context.Context, endpoint string) ([]model.Gateway, error) {
	query := g.baseQuery().Where(goqu.Ex{"endpoint": endpoint})
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return nil, err
	}

	var gateways []model.Gateway
	if err := g.Store.SelectContext(ctx, &gateways, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return nil, err
	}

	return gateways, nil
}

func (g *Gateway) GetByLogLevel(ctx context.Context, logLevel string) ([]model.Gateway, error) {
	query := g.baseQuery().Where(goqu.Ex{"log_level": logLevel})
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return nil, err
	}

	var gateways []model.Gateway
	if err := g.Store.SelectContext(ctx, &gateways, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return nil, err
	}

	return gateways, nil
}

func (g *Gateway) GetByStage(ctx context.Context, stage string) ([]model.Gateway, error) {
	query := g.baseQuery().Where(goqu.Ex{"stage": stage})
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return nil, err
	}

	var gateways []model.Gateway
	if err := g.Store.SelectContext(ctx, &gateways, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return nil, err
	}

	return gateways, nil
}

func (g *Gateway) GetByName(ctx context.Context, name string) ([]model.Gateway, error) {
	query := g.baseQuery().Where(goqu.Ex{"name": name})
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return nil, err
	}

	var gateways []model.Gateway
	if err := g.Store.SelectContext(ctx, &gateways, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return nil, err
	}

	return gateways, nil
}

func (g *Gateway) GetByApplicationID(ctx context.Context, applicationID string) ([]model.Gateway, error) {
	query := g.baseQuery().Where(goqu.Ex{"application_id": applicationID})
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return nil, err
	}

	var gateways []model.Gateway
	if err := g.Store.SelectContext(ctx, &gateways, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return nil, err
	}

	return gateways, nil
}
