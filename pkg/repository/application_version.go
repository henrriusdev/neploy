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

type ApplicationVersion struct {
	Base[model.ApplicationVersion]
}

func NewApplicationVersion(db store.Queryable) *ApplicationVersion {
	return &ApplicationVersion{Base[model.ApplicationVersion]{Store: db, Table: "application_versions"}}
}

func (a *ApplicationVersion) Insert(ctx context.Context, version model.ApplicationVersion) error {
	_, err := a.UpsertOneDoNothing(ctx, version, "application_id", "version_tag")
	return err
}

func (a *ApplicationVersion) Delete(ctx context.Context, id string) error {
	query := filters.ApplyUpdateFilters(
		a.BaseQueryUpdate().Set(goqu.Record{"deleted_at": goqu.L("CURRENT_TIMESTAMP")}),
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

func (a *ApplicationVersion) Exists(ctx context.Context, appID, tag string) (bool, error) {
	row, err := a.GetOne(ctx, filters.IsSelectFilter("application_id", appID), filters.IsSelectFilter("version_tag", tag))
	return row.ID != "", err
}

func (a *ApplicationVersion) ExistsByName(ctx context.Context, name, tag string) (bool, error) {
	query := a.baseQuery("v").
		Select(goqu.I("v.*")).
		LeftJoin(
			goqu.T("gateways").As("g"),
			goqu.On(goqu.I("g.application_id").Eq(goqu.I("v.application_id"))),
		).Where(goqu.I("g.path").Eq("/" + name)).Where(goqu.I("v.version_tag").Eq(tag))

	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return false, err
	}

	var row model.ApplicationVersion
	if err := a.Store.GetContext(ctx, &row, q, args...); err != nil {
		logger.Error("error executing exists by name %s: %v", name, err)
		return false, err
	}

	common.AttachSQLToTrace(ctx, q)
	return row.ID != "", err
}
