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

type VisitorInfo struct {
	Base[model.VisitorInfo]
}

func NewVisitor(db store.Queryable) *VisitorInfo {
	return &VisitorInfo{Base[model.VisitorInfo]{Store: db, Table: "visitor_info"}}
}

func (v *VisitorInfo) GetByID(ctx context.Context, id string) (model.VisitorInfo, error) {
	query := v.baseQuery().Where(goqu.Ex{"id": id})
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return model.VisitorInfo{}, err
	}

	var visitorInfo model.VisitorInfo
	if err := v.Store.GetContext(ctx, &visitorInfo, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return model.VisitorInfo{}, err
	}

	common.AttachSQLToTrace(ctx, q)
	return visitorInfo, nil
}

func (v *VisitorInfo) Insert(ctx context.Context, visitorInfo model.VisitorInfo) error {
	query := v.BaseQueryInsert().Rows(visitorInfo)
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

func (v *VisitorInfo) Update(ctx context.Context, visitorInfo model.VisitorInfo) error {
	query := filters.ApplyUpdateFilters(v.BaseQueryUpdate().Set(visitorInfo), filters.IsUpdateFilter("id", visitorInfo.ID))
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

func (v *VisitorInfo) Delete(ctx context.Context, id string) error {
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

func (v *VisitorInfo) GetAll(ctx context.Context) ([]model.VisitorInfo, error) {
	query := v.baseQuery()
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return nil, err
	}

	var visitorInfos []model.VisitorInfo
	if err := v.Store.SelectContext(ctx, &visitorInfos, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return nil, err
	}

	common.AttachSQLToTrace(ctx, q)
	return visitorInfos, nil
}
