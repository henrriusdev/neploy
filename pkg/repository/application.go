package repository

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository/filters"
	"neploy.dev/pkg/store"
)

type Application interface {
	Insert(ctx context.Context, application model.Application) (string, error)
	Update(ctx context.Context, application model.Application) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (model.Application, error)
	GetAll(ctx context.Context) ([]model.Application, error)
}

type application[T any] struct {
	Base[T]
}

func NewApplication(db store.Queryable) Application {
	return &application[model.Application]{Base[model.Application]{Store: db, Table: "applications"}}
}

func (a *application[T]) Insert(ctx context.Context, application model.Application) (string, error) {
	var id string
	query := a.BaseQueryInsert().Rows(application).Returning("id")
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building insert query: %v", err)
		return "", err
	}

	if err := a.Store.QueryRowxContext(ctx, q, args...).Scan(&id); err != nil {
		logger.Error("error executing insert query: %v", err)
		return "", err
	}

	return id, nil
}

func (a *application[T]) Update(ctx context.Context, application model.Application) error {
	query := filters.ApplyUpdateFilters(a.BaseQueryUpdate().Set(application), filters.IsUpdateFilter("id", application.ID))
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

func (a *application[T]) Delete(ctx context.Context, id string) error {
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

func (a *application[T]) GetByID(ctx context.Context, id string) (model.Application, error) {
	query := a.baseQuery().Where(goqu.Ex{"id": id})
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return model.Application{}, err
	}

	var application model.Application
	if err := a.Store.QueryRowxContext(ctx, q, args...).StructScan(&application); err != nil {
		logger.Error("error executing select query: %v", err)
		return model.Application{}, err
	}

	return application, nil
}

func (a *application[T]) GetAll(ctx context.Context) ([]model.Application, error) {
	query := a.baseQuery()
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return nil, err
	}

	var applications []model.Application
	if err := a.Store.SelectContext(ctx, &applications, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return nil, err
	}

	return applications, nil
}
