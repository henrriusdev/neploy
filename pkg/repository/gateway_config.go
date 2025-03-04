package repository

import (
	"context"
	"database/sql"
	"errors"
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository/filters"
	"neploy.dev/pkg/store"
)

type GatewayConfig interface {
	Upsert(ctx context.Context, gateway model.GatewayConfig) (model.GatewayConfig, error)
	Get(ctx context.Context) (model.GatewayConfig, error)
}

type gatewayConfig[T any] struct {
	Base[T]
}

func NewGatewayConfig(db store.Queryable) GatewayConfig {
	return &gatewayConfig[model.GatewayConfig]{Base[model.GatewayConfig]{Store: db, Table: "gateway_config"}}
}

func (g *gatewayConfig[T]) Upsert(ctx context.Context, gateway model.GatewayConfig) (model.GatewayConfig, error) {
	conf, err := g.Get(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Error("error getting actual gateway: %v", err)
		return model.GatewayConfig{}, err
	}

	if conf.ID != "" {
		query := g.BaseQueryUpdate().Set(gateway).Returning("*")
		q, args, err := query.ToSQL()
		if err != nil {
			logger.Error("error building upsert query: %v", err)
			return model.GatewayConfig{}, err
		}

		if err := g.Store.QueryRowxContext(ctx, q, args...).StructScan(&conf); err != nil {
			logger.Error("error running upsert query: %v", err)
			return model.GatewayConfig{}, err
		}

		return conf, nil
	}

	query := g.BaseQueryInsert().Rows(gateway).Returning("*")
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building insert query: %v", err)
		return model.GatewayConfig{}, err
	}

	if err := g.Store.QueryRowxContext(ctx, q, args...).StructScan(&conf); err != nil {
		logger.Error("error executing insert query: %v", err)
		return model.GatewayConfig{}, err
	}

	return conf, nil
}

func (g *gatewayConfig[T]) Get(ctx context.Context) (conf model.GatewayConfig, err error) {
	query := filters.ApplyFilters(g.baseQuery(), filters.LimitOffsetFilter(1, 0))
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building get query: %v", err)
		return model.GatewayConfig{}, err
	}

	if err = g.Store.GetContext(ctx, &conf, q, args...); err != nil {
		logger.Error("error running get config query: %v", err)
		return model.GatewayConfig{}, err
	}

	return
}

func (g *gatewayConfig[T]) createDefault(ctx context.Context) (conf model.GatewayConfig, err error) {
	conf, err = g.InsertOne(ctx, model.GatewayConfig{
		DefaultVersioningType: "headers",
		DefaultVersion:        "latest",
		LoadBalancer:          false,
	})

	return
}
