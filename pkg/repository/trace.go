package repository

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository/filters"
	"neploy.dev/pkg/store"
)

type Trace struct {
	Base[model.Trace]
}

func NewTrace(db store.Queryable) *Trace {
	return &Trace{Base[model.Trace]{Store: db, Table: "traces"}}
}

func (t *Trace) Insert(ctx context.Context, trace model.Trace) error {
	query := t.BaseQueryInsert().Rows(trace)
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building insert query: %v", err)
		return err
	}

	if _, err := t.Store.ExecContext(ctx, q, args...); err != nil {
		logger.Error("error executing insert query: %v", err)
		return err
	}

	return nil
}

func (t *Trace) Update(ctx context.Context, trace model.Trace) error {
	query := filters.ApplyUpdateFilters(t.BaseQueryUpdate().Set(trace), filters.IsUpdateFilter("id", trace.ID))
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building update query: %v", err)
		return err
	}

	if _, err := t.Store.ExecContext(ctx, q, args...); err != nil {
		logger.Error("error executing update query: %v", err)
		return err
	}

	return nil
}

func (t *Trace) Delete(ctx context.Context, id string) error {
	query := filters.ApplyUpdateFilters(
		t.BaseQueryUpdate().
			Set(goqu.Record{"deleted_at": goqu.L("CURRENT_TIMESTAMP")}),
		filters.IsUpdateFilter("id", id),
	)

	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building delete query: %v", err)
		return err
	}

	if _, err := t.Store.ExecContext(ctx, q, args...); err != nil {
		logger.Error("error executing delete query: %v", err)
		return err
	}

	return nil
}

func (t *Trace) GetByID(ctx context.Context, id string) (model.Trace, error) {
	query := t.baseQuery().Where(goqu.Ex{"id": id})
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return model.Trace{}, err
	}

	var trace model.Trace
	if err := t.Store.GetContext(ctx, &trace, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return model.Trace{}, err
	}

	return trace, nil
}

func (t *Trace) GetAll(ctx context.Context) ([]model.Trace, error) {
	query := t.baseQuery()
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return nil, err
	}

	var traces []model.Trace
	if err := t.Store.SelectContext(ctx, &traces, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return nil, err
	}

	return traces, nil
}

func (t *Trace) GetByUserID(ctx context.Context, userID string) ([]model.Trace, error) {
	query := t.baseQuery().Where(goqu.Ex{"user_id": userID})
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return nil, err
	}

	var traces []model.Trace
	if err := t.Store.SelectContext(ctx, &traces, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return nil, err
	}

	return traces, nil
}

func (t *Trace) GetByType(ctx context.Context, traceType string) ([]model.Trace, error) {
	query := t.baseQuery().Where(goqu.Ex{"type": traceType})
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return nil, err
	}

	var traces []model.Trace
	if err := t.Store.SelectContext(ctx, &traces, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return nil, err
	}

	return traces, nil
}

func (t *Trace) GetByAction(ctx context.Context, action string) ([]model.Trace, error) {
	query := t.baseQuery().Where(goqu.Ex{"action": action})
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return nil, err
	}

	var traces []model.Trace
	if err := t.Store.SelectContext(ctx, &traces, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return nil, err
	}

	return traces, nil
}

func (t *Trace) GetByActionTimestamp(ctx context.Context, timestamp model.Date) ([]model.Trace, error) {
	query := t.baseQuery().Where(goqu.Ex{"timestamp": timestamp})
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return nil, err
	}

	var traces []model.Trace
	if err := t.Store.SelectContext(ctx, &traces, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return nil, err
	}

	return traces, nil
}
