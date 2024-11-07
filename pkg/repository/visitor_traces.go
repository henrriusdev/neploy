package repository

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/rs/zerolog/log"
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
		log.Err(err).Msg("error building insert query")
		return err
	}

	if _, err := v.Store.ExecContext(ctx, q, args...); err != nil {
		log.Err(err).Msg("error executing insert query")
		return err
	}

	return nil
}

func (v *visitorTraces[T]) Update(ctx context.Context, visitorTraces model.VisitorTrace) error {
	query := filters.ApplyUpdateFilters(v.BaseQueryUpdate().Set(visitorTraces), filters.IsUpdateFilter("id", visitorTraces.ID))
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building update query")
		return err
	}

	if _, err := v.Store.ExecContext(ctx, q, args...); err != nil {
		log.Err(err).Msg("error executing update query")
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
		log.Err(err).Msg("error building delete query")
		return err
	}

	if _, err := v.Store.ExecContext(ctx, q, args...); err != nil {
		log.Err(err).Msg("error executing delete query")
		return err
	}

	return nil
}

func (v *visitorTraces[T]) GetByID(ctx context.Context, id string) (model.VisitorTrace, error) {
	query := v.baseQuery().Where(goqu.C("id").Eq(id))
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building select query")
		return model.VisitorTrace{}, err
	}

	var visitorTraces model.VisitorTrace
	if err := v.Store.GetContext(ctx, &visitorTraces, q, args...); err != nil {
		log.Err(err).Msg("error executing select query")
		return model.VisitorTrace{}, err
	}

	return visitorTraces, nil
}

func (v *visitorTraces[T]) GetAll(ctx context.Context) ([]model.VisitorTrace, error) {
	query := v.baseQuery()
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building select query")
		return nil, err
	}

	var visitorTraces []model.VisitorTrace
	if err := v.Store.SelectContext(ctx, &visitorTraces, q, args...); err != nil {
		log.Err(err).Msg("error executing select query")
		return nil, err
	}

	return visitorTraces, nil
}

func (v *visitorTraces[T]) GetByVisitorID(ctx context.Context, visitorID string) ([]model.VisitorTrace, error) {
	query := v.baseQuery().Where(goqu.C("visitor_id").Eq(visitorID))
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building select query")
		return nil, err
	}

	var visitorTraces []model.VisitorTrace
	if err := v.Store.SelectContext(ctx, &visitorTraces, q, args...); err != nil {
		log.Err(err).Msg("error executing select query")
		return nil, err
	}

	return visitorTraces, nil
}
