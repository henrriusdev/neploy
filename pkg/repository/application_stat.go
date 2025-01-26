package repository

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository/filters"
	"neploy.dev/pkg/store"
)

type ApplicationStat interface {
	Insert(ctx context.Context, applicationStat model.ApplicationStat) error
	Update(ctx context.Context, applicationStat model.ApplicationStat) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (model.ApplicationStat, error)
	GetByApplicationID(ctx context.Context, applicationID string) ([]model.ApplicationStat, error)
	GetByEnvironmentID(ctx context.Context, environmentID string) ([]model.ApplicationStat, error)
	GetByApplicationIDAndEnvironmentID(ctx context.Context, applicationID, environmentID string) (model.ApplicationStat, error)
	GetByDate(ctx context.Context, date time.Time) ([]model.ApplicationStat, error)
	GetAll(ctx context.Context) ([]model.ApplicationStat, error)
	GetUniqueVisitors(ctx context.Context, applicationID, environmentID string) (int, error)
	GetDataTransfered(ctx context.Context, applicationID, environmentID string) (int, error)
	GetRequests(ctx context.Context, applicationID, environmentID string) (int, error)
	GetAverageResponseTime(ctx context.Context, applicationID, environmentID string) (int, error)
	GetErrorRate(ctx context.Context, applicationID, environmentID string) (int, error)
}

type applicationStat[T any] struct {
	Base[T]
}

func NewApplicationStat(db store.Queryable) ApplicationStat {
	return &applicationStat[model.ApplicationStat]{Base[model.ApplicationStat]{Store: db, Table: "application_stats"}}
}

func (a *applicationStat[T]) Insert(ctx context.Context, applicationStat model.ApplicationStat) error {
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

	return nil
}

func (a *applicationStat[T]) Update(ctx context.Context, applicationStat model.ApplicationStat) error {
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

	return nil
}

func (a *applicationStat[T]) Delete(ctx context.Context, id string) error {
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

func (a *applicationStat[T]) GetByID(ctx context.Context, id string) (model.ApplicationStat, error) {
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

	return applicationStat, nil
}

func (a *applicationStat[T]) GetByApplicationID(ctx context.Context, applicationID string) ([]model.ApplicationStat, error) {
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

	return applicationStats, nil
}

func (a *applicationStat[T]) GetByEnvironmentID(ctx context.Context, environmentID string) ([]model.ApplicationStat, error) {
	query := a.baseQuery().Where(goqu.Ex{"environment_id": environmentID})
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

	return applicationStats, nil
}

func (a *applicationStat[T]) GetByDate(ctx context.Context, date time.Time) ([]model.ApplicationStat, error) {
	query := a.baseQuery().Where(goqu.Ex{"date": date})
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

	return applicationStats, nil
}

func (a *applicationStat[T]) GetAll(ctx context.Context) ([]model.ApplicationStat, error) {
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

	return applicationStats, nil
}

func (a *applicationStat[T]) GetUniqueVisitors(ctx context.Context, applicationID, environmentID string) (int, error) {
	query := a.baseQuery().Where(goqu.Ex{"application_id": applicationID, "environment_id": environmentID}).Select(goqu.COUNT("DISTINCT(visitor_id)"))
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return 0, err
	}

	var count int
	if err := a.Store.GetContext(ctx, &count, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return 0, err
	}

	return count, nil
}

func (a *applicationStat[T]) GetDataTransfered(ctx context.Context, applicationID, environmentID string) (int, error) {
	query := a.baseQuery().Where(goqu.Ex{"application_id": applicationID, "environment_id": environmentID}).Select(goqu.SUM("data_transfered"))
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return 0, err
	}

	var sum int
	if err := a.Store.GetContext(ctx, &sum, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return 0, err
	}

	return sum, nil
}

func (a *applicationStat[T]) GetRequests(ctx context.Context, applicationID, environmentID string) (int, error) {
	query := a.baseQuery().Where(goqu.Ex{"application_id": applicationID, "environment_id": environmentID}).Select(goqu.SUM("requests"))
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return 0, err
	}

	var sum int
	if err := a.Store.GetContext(ctx, &sum, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return 0, err
	}

	return sum, nil
}

func (a *applicationStat[T]) GetAverageResponseTime(ctx context.Context, applicationID, environmentID string) (int, error) {
	query := a.baseQuery().Where(goqu.Ex{"application_id": applicationID, "environment_id": environmentID}).Select(goqu.AVG("response_time"))
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return 0, err
	}

	var avg int
	if err := a.Store.GetContext(ctx, &avg, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return 0, err
	}

	return avg, nil
}

func (a *applicationStat[T]) GetErrorRate(ctx context.Context, applicationID, environmentID string) (int, error) {
	query := a.baseQuery().Where(goqu.Ex{"application_id": applicationID, "environment_id": environmentID}).Select(goqu.AVG("error_rate"))
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return 0, err
	}

	var avg int
	if err := a.Store.GetContext(ctx, &avg, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return 0, err
	}

	return avg, nil
}

func (a *applicationStat[T]) GetByApplicationIDAndEnvironmentID(ctx context.Context, applicationID, environmentID string) (model.ApplicationStat, error) {
	query := a.baseQuery().Where(goqu.Ex{"application_id": applicationID, "environment_id": environmentID})
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

	return applicationStat, nil
}
