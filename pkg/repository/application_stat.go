package repository

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	"neploy.dev/pkg/common"
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository/filters"
	"neploy.dev/pkg/store"
)

type ApplicationStat struct {
	Base[model.ApplicationStat]
}

func NewApplicationStat(db store.Queryable) *ApplicationStat {
	return &ApplicationStat{Base[model.ApplicationStat]{Store: db, Table: "application_stats"}}
}

func (a *ApplicationStat) Insert(ctx context.Context, applicationStat model.ApplicationStat) error {
	query := a.BaseQueryInsert().Rows(applicationStat)
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building insert query: %v", err)
		return err
	}

	if _, err := a.Store.ExecContext(ctx, q, args...); err != nil {
		logger.Error("error executing insert query: %v", err)
		return err
	}

	common.AttachSQLToTrace(ctx, q)
	return nil
}

func (a *ApplicationStat) Update(ctx context.Context, applicationStat model.ApplicationStat) error {
	query := filters.ApplyUpdateFilters(a.BaseQueryUpdate().Set(applicationStat), filters.IsUpdateFilter("id", applicationStat.ID))
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building update query: %v", err)
		return err
	}

	if _, err := a.Store.ExecContext(ctx, q, args...); err != nil {
		logger.Error("error executing update query: %v", err)
		return err
	}

	common.AttachSQLToTrace(ctx, q)
	return nil
}

func (a *ApplicationStat) Delete(ctx context.Context, id string) error {
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

	common.AttachSQLToTrace(ctx, q)
	return nil
}

func (a *ApplicationStat) GetByID(ctx context.Context, id string) (model.ApplicationStat, error) {
	query := a.baseQuery().Where(goqu.Ex{"id": id})
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return model.ApplicationStat{}, err
	}

	var applicationStat model.ApplicationStat
	if err := a.Store.GetContext(ctx, &applicationStat, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return model.ApplicationStat{}, err
	}

	common.AttachSQLToTrace(ctx, q)
	return applicationStat, nil
}

func (a *ApplicationStat) GetByApplicationID(ctx context.Context, applicationID string) ([]model.ApplicationStat, error) {
	query := a.baseQuery().Where(goqu.Ex{"application_id": applicationID})
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return nil, err
	}

	var applicationStats []model.ApplicationStat
	if err := a.Store.SelectContext(ctx, &applicationStats, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return nil, err
	}

	common.AttachSQLToTrace(ctx, q)
	return applicationStats, nil
}

func (a *ApplicationStat) GetAll(ctx context.Context) ([]model.ApplicationStat, error) {
	query := a.baseQuery()
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return nil, err
	}

	var applicationStats []model.ApplicationStat
	if err := a.Store.SelectContext(ctx, &applicationStats, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return nil, err
	}

	common.AttachSQLToTrace(ctx, q)
	return applicationStats, nil
}

func (a *ApplicationStat) GetHourlyRequests(ctx context.Context) ([]model.RequestStat, error) {
	query := goqu.
		From("application_stats").
		Select(
			goqu.L("to_char(date AT TIME ZONE 'UTC', 'YYYY-MM-DD HH24:00')").As("hour"),
			goqu.SUM("requests").As("successful"),
			goqu.SUM("errors").As("errors"),
		).
		GroupBy(goqu.L("hour")).
		Order(goqu.L("hour").Asc()).
		Limit(24)

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, err
	}

	var stats []model.RequestStat
	if err := a.Store.SelectContext(ctx, &stats, sql, args...); err != nil {
		return nil, err
	}

	return stats, nil
}
