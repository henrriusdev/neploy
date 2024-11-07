package repository

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/rs/zerolog/log"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository/filters"
	"neploy.dev/pkg/store"
)

type RefreshToken interface {
	Insert(ctx context.Context, refreshToken model.RefreshToken) error
	Update(ctx context.Context, refreshToken model.RefreshToken) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (model.RefreshToken, error)
	GetByUserID(ctx context.Context, userID string) (model.RefreshToken, error)
	GetAll(ctx context.Context) ([]model.RefreshToken, error)
}

type refreshToken[T any] struct {
	Base[T]
}

func NewRefreshToken(db store.Queryable) RefreshToken {
	return &refreshToken[model.RefreshToken]{Base[model.RefreshToken]{Store: db, Table: "refresh_tokens"}}
}

func (r *refreshToken[T]) Insert(ctx context.Context, refreshToken model.RefreshToken) error {
	query := r.BaseQueryInsert().Rows(refreshToken)
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building insert query")
		return err
	}

	if _, err := r.Store.ExecContext(ctx, q, args...); err != nil {
		log.Err(err).Msg("error executing insert query")
		return err
	}

	return nil
}

func (r *refreshToken[T]) Update(ctx context.Context, refreshToken model.RefreshToken) error {
	query := filters.ApplyUpdateFilters(r.BaseQueryUpdate().Set(refreshToken), filters.IsUpdateFilter("id", refreshToken.ID))
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building update query")
		return err
	}

	if _, err := r.Store.ExecContext(ctx, q, args...); err != nil {
		log.Err(err).Msg("error executing update query")
		return err
	}

	return nil
}

func (r *refreshToken[T]) Delete(ctx context.Context, id string) error {
	query := filters.ApplyUpdateFilters(
		r.BaseQueryUpdate().
			Set(goqu.Record{"deleted_at": goqu.L("CURRENT_TIMESTAMP")}),
		filters.IsUpdateFilter("id", id),
	)

	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building delete query")
		return err
	}

	if _, err := r.Store.ExecContext(ctx, q, args...); err != nil {
		log.Err(err).Msg("error executing delete query")
		return err
	}

	return nil
}

func (r *refreshToken[T]) GetByID(ctx context.Context, id string) (model.RefreshToken, error) {
	query := r.baseQuery().Where(goqu.Ex{"id": id})
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building select query")
		return model.RefreshToken{}, err
	}

	var refreshToken model.RefreshToken
	if err := r.Store.GetContext(ctx, &refreshToken, q, args...); err != nil {
		log.Err(err).Msg("error executing select query")
		return model.RefreshToken{}, err
	}

	return refreshToken, nil
}

func (r *refreshToken[T]) GetByUserID(ctx context.Context, userID string) (model.RefreshToken, error) {
	query := r.baseQuery().Where(goqu.Ex{"user_id": userID})
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building select query")
		return model.RefreshToken{}, err
	}

	var refreshToken model.RefreshToken
	if err := r.Store.GetContext(ctx, &refreshToken, q, args...); err != nil {
		log.Err(err).Msg("error executing select query")
		return model.RefreshToken{}, err
	}

	return refreshToken, nil
}

func (r *refreshToken[T]) GetAll(ctx context.Context) ([]model.RefreshToken, error) {
	query := r.baseQuery()
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building select query")
		return nil, err
	}

	var refreshTokens []model.RefreshToken
	if err := r.Store.SelectContext(ctx, &refreshTokens, q, args...); err != nil {
		log.Err(err).Msg("error executing select query")
		return nil, err
	}

	return refreshTokens, nil
}
