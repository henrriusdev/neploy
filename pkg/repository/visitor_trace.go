package repository

import (
	"context"
	"neploy.dev/pkg/common"

	"github.com/doug-martin/goqu/v9"
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository/filters"
	"neploy.dev/pkg/store"
)

type VisitorTrace struct {
	Base[model.VisitorTrace]
}

func NewVisitorTrace(db store.Queryable) *VisitorTrace {
	return &VisitorTrace{Base[model.VisitorTrace]{Store: db, Table: "visitor_traces"}}
}

func (v *VisitorTrace) Insert(ctx context.Context, visitorTraces model.VisitorTrace) error {
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

	common.AttachSQLToTrace(ctx, q)
	return nil
}

func (v *VisitorTrace) Update(ctx context.Context, visitorTraces model.VisitorTrace) error {
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

	common.AttachSQLToTrace(ctx, q)
	return nil
}

func (v *VisitorTrace) Delete(ctx context.Context, id string) error {
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

	common.AttachSQLToTrace(ctx, q)
	return nil
}

func (v *VisitorTrace) GetByID(ctx context.Context, id string) (model.VisitorTrace, error) {
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

	common.AttachSQLToTrace(ctx, q)
	return visitorTraces, nil
}

func (v *VisitorTrace) GetAll(ctx context.Context) ([]model.VisitorTrace, error) {
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

	common.AttachSQLToTrace(ctx, q)
	return visitorTraces, nil
}

func (v *VisitorTrace) GetByVisitorID(ctx context.Context, visitorID string) ([]model.VisitorTrace, error) {
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

	common.AttachSQLToTrace(ctx, q)
	return visitorTraces, nil
}

func (v *VisitorTrace) Create(ctx context.Context, visitorTrace model.VisitorTrace, resolvedVersion string) (model.VisitorTrace, error) {
	exists, err := NewApplicationVersion(v.Store).ExistsByName(ctx, visitorTrace.ApplicationID, resolvedVersion)
	if err != nil || !exists {
		logger.Error("error checking application version existence %v", err)
		return model.VisitorTrace{}, err
	}

	gateway, err := NewGateway(v.Store).GetByPath(ctx, visitorTrace.ApplicationID)
	if err != nil {
		logger.Error("error getting application by name %v", err)
		return model.VisitorTrace{}, err
	}

	visitorTrace.ApplicationID = gateway.ApplicationID
	trace, err := v.InsertOne(ctx, visitorTrace)
	if err != nil {
		logger.Error("error inserting visitor trace %v", err)
	}

	return trace, err
}

func (v *VisitorTrace) GetTraces(ctx context.Context) ([]model.VisitorStat, error) {
	query := v.baseQuery().Select(
		goqu.COUNT(goqu.DISTINCT(goqu.C("id"))).As("amount"),
		goqu.L("DATE(visit_timestamp)").As("date"),
	).
		GroupBy(goqu.L("DATE(visit_timestamp)")).
		Order(goqu.L("DATE(visit_timestamp)").Asc()).
		Limit(1000)

	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return nil, err
	}

	var visitorTraces []model.VisitorStat
	if err := v.Store.SelectContext(ctx, &visitorTraces, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return nil, err
	}

	common.AttachSQLToTrace(ctx, q)
	return visitorTraces, nil
}
