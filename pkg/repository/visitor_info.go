package repository

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository/filters"
	"neploy.dev/pkg/store"
)

type VisitorInfo interface {
	Insert(ctx context.Context, visitorInfo model.VisitorInfo) error
	Update(ctx context.Context, visitorInfo model.VisitorInfo) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (model.VisitorInfo, error)
	GetAll(ctx context.Context) ([]model.VisitorInfo, error)
}

type visitorInfo[T any] struct {
	Base[T]
}

func NewVisitor(db store.Queryable) VisitorInfo {
	return &visitorInfo[model.VisitorInfo]{Base[model.VisitorInfo]{Store: db, Table: "visitor_info"}}
}

func (v *visitorInfo[T]) GetByID(ctx context.Context, id string) (model.VisitorInfo, error) {
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

	return visitorInfo, nil
}

func (v *visitorInfo[T]) Insert(ctx context.Context, visitorInfo model.VisitorInfo) error {
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

	return nil
}

func (v *visitorInfo[T]) Update(ctx context.Context, visitorInfo model.VisitorInfo) error {
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

	return nil
}

func (v *visitorInfo[T]) Delete(ctx context.Context, id string) error {
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

func (v *visitorInfo[T]) GetAll(ctx context.Context) ([]model.VisitorInfo, error) {
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

	return visitorInfos, nil
}
