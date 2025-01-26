package repository

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository/filters"
	"neploy.dev/pkg/store"
)

type VisitorTrace interface {
	Insert(ctx context.Context, visitorTraces model.VisitorTrace) error
	Update(ctx context.Context, visitorTraces model.VisitorTrace) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (model.VisitorTrace, error)
	GetAll(ctx context.Context) ([]model.VisitorTrace, error)
	GetByVisitorID(ctx context.Context, visitorID string) ([]model.VisitorTrace, error)
}

type visitorTraces[T any] struct {
	Base[T]
}

func NewVisitorTrace(db store.Queryable) VisitorTrace {
	return &visitorTraces[model.VisitorTrace]{Base[model.VisitorTrace]{Store: db, Table: "visitor_traces"}}
}

func (v *visitorTraces[T]) Insert(ctx context.Context, visitorTraces model.VisitorTrace) error {
	query := v.BaseQueryInsert().Rows(visitorTraces)
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building insert query: %v", err)
		return err
	}

	if _, err := v.Store.ExecContext(ctx, q, args...); err != nil {
		logger.Error("error executing insert query: %v", err)
		return err
	}

	return nil
}

func (v *visitorTraces[T]) Update(ctx context.Context, visitorTraces model.VisitorTrace) error {
	query := filters.ApplyUpdateFilters(v.BaseQueryUpdate().Set(visitorTraces), filters.IsUpdateFilter("id", visitorTraces.ID))
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building update query: %v", err)
		return err
	}

	if _, err := v.Store.ExecContext(ctx, q, args...); err != nil {
		logger.Error("error executing update query: %v", err)
		return err
	}

	return nil
}

func (v *visitorTraces[T]) Delete(ctx context.Context, id string) error {
	query := filters.ApplyUpdateFilters(
		v.BaseQueryUpdate().
			Set(goqu.Record{"deleted_at": goqu.L("CURRENT_TIMESTAMP")}),
		filters.IsUpdateFilter("id", id),
	)

	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building delete query: %v", err)
		return err
	}

	if _, err := v.Store.ExecContext(ctx, q, args...); err != nil {
		logger.Error("error executing delete query: %v", err)
		return err
	}

	return nil
}

func (v *visitorTraces[T]) GetByID(ctx context.Context, id string) (model.VisitorTrace, error) {
	query := v.baseQuery().Where(goqu.C("id").Eq(id))
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return model.VisitorTrace{}, err
	}

	var visitorTraces model.VisitorTrace
	if err := v.Store.GetContext(ctx, &visitorTraces, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return model.VisitorTrace{}, err
	}

	return visitorTraces, nil
}

func (v *visitorTraces[T]) GetAll(ctx context.Context) ([]model.VisitorTrace, error) {
	query := v.baseQuery()
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return nil, err
	}

	var visitorTraces []model.VisitorTrace
	if err := v.Store.SelectContext(ctx, &visitorTraces, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return nil, err
	}

	return visitorTraces, nil
}

func (v *visitorTraces[T]) GetByVisitorID(ctx context.Context, visitorID string) ([]model.VisitorTrace, error) {
	query := v.baseQuery().Where(goqu.C("visitor_id").Eq(visitorID))
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return nil, err
	}

	var visitorTraces []model.VisitorTrace
	if err := v.Store.SelectContext(ctx, &visitorTraces, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return nil, err
	}

	return visitorTraces, nil
}
