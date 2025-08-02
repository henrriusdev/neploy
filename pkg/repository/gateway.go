package repository

import (
	"context"
	"errors"

	"neploy.dev/pkg/common"

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
	if _, err := g.GetByPath(ctx, gateway.Path); err == nil {
		return errors.New("path already exists")
	}

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

	common.AttachSQLToTrace(ctx, q)
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

	common.AttachSQLToTrace(ctx, q)
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

	common.AttachSQLToTrace(ctx, q)
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

	common.AttachSQLToTrace(ctx, q)
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

	common.AttachSQLToTrace(ctx, q)
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

	common.AttachSQLToTrace(ctx, q)
	return gateways, nil
}

func (g *Gateway) GetByPath(ctx context.Context, path string) (model.Gateway, error) {
	query := g.baseQuery().Where(goqu.Ex{"path": "/" + path})
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

	common.AttachSQLToTrace(ctx, q)
	return gateway, nil
}
