package repository

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository/filters"
	"neploy.dev/pkg/store"
)

type RefreshToken struct {
	Base[model.RefreshToken]
}

func NewRefreshToken(db store.Queryable) *RefreshToken {
	return &RefreshToken{Base[model.RefreshToken]{Store: db, Table: "refresh_tokens"}}
}

func (r *RefreshToken) Insert(ctx context.Context, refreshToken model.RefreshToken) error {
	query := r.BaseQueryInsert().Rows(refreshToken)
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building insert query: %v", err)
		return err
	}

	if _, err := r.Store.ExecContext(ctx, q, args...); err != nil {
		logger.Error("error executing insert query: %v", err)
		return err
	}

	return nil
}

func (r *RefreshToken) Update(ctx context.Context, refreshToken model.RefreshToken) error {
	query := filters.ApplyUpdateFilters(r.BaseQueryUpdate().Set(refreshToken), filters.IsUpdateFilter("id", refreshToken.ID))
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building update query: %v", err)
		return err
	}

	if _, err := r.Store.ExecContext(ctx, q, args...); err != nil {
		logger.Error("error executing update query: %v", err)
		return err
	}

	return nil
}

func (r *RefreshToken) Delete(ctx context.Context, id string) error {
	query := filters.ApplyUpdateFilters(
		r.BaseQueryUpdate().
			Set(goqu.Record{"deleted_at": goqu.L("CURRENT_TIMESTAMP")}),
		filters.IsUpdateFilter("id", id),
	)

	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building delete query: %v", err)
		return err
	}

	if _, err := r.Store.ExecContext(ctx, q, args...); err != nil {
		logger.Error("error executing delete query: %v", err)
		return err
	}

	return nil
}

func (r *RefreshToken) GetByID(ctx context.Context, id string) (model.RefreshToken, error) {
	query := r.baseQuery().Where(goqu.Ex{"id": id})
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return model.RefreshToken{}, err
	}

	var refreshToken model.RefreshToken
	if err := r.Store.GetContext(ctx, &refreshToken, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return model.RefreshToken{}, err
	}

	return refreshToken, nil
}

func (r *RefreshToken) GetByUserID(ctx context.Context, userID string) (model.RefreshToken, error) {
	query := r.baseQuery().Where(goqu.Ex{"user_id": userID})
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return model.RefreshToken{}, err
	}

	var refreshToken model.RefreshToken
	if err := r.Store.GetContext(ctx, &refreshToken, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return model.RefreshToken{}, err
	}

	return refreshToken, nil
}

func (r *RefreshToken) GetAll(ctx context.Context) ([]model.RefreshToken, error) {
	query := r.baseQuery()
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return nil, err
	}

	var refreshTokens []model.RefreshToken
	if err := r.Store.SelectContext(ctx, &refreshTokens, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return nil, err
	}

	return refreshTokens, nil
}
