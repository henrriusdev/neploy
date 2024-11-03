package repository

import (
	"context"

	"github.com/rs/zerolog/log"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository/filters"
	"neploy.dev/pkg/store"
)

type UserOauth interface {
	Insert(ctx context.Context, user model.UserOAuth) error
	GetByOAuthID(ctx context.Context, oauthID string) (model.UserOAuth, error)
	GetByUserID(ctx context.Context, userID string) (model.UserOAuth, error)
}

type userOauth[T any] struct {
	Base[T]
}

func NewUserOauth(db store.Queryable) UserOauth {
	return &userOauth[model.UserOAuth]{Base[model.UserOAuth]{Store: db, Table: "user_oauth"}}
}

func (u *userOauth[T]) Insert(ctx context.Context, user model.UserOAuth) error {
	query := u.BaseQueryInsert().Rows(user)
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building query")
		return err
	}

	_, err = u.Store.ExecContext(ctx, q, args...)
	if err != nil {
		log.Err(err).Msg("error inserting user oauth")
		return err
	}

	return nil
}

func (u *userOauth[T]) GetByOAuthID(ctx context.Context, oauthID string) (model.UserOAuth, error) {
	query := filters.ApplyFilters(u.baseQuery(), filters.IsSelectFilter("oauth_id", oauthID))
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building query")
		return model.UserOAuth{}, err
	}

	var user model.UserOAuth
	err = u.Store.GetContext(ctx, &user, q, args...)
	if err != nil {
		log.Err(err).Msg("error getting user oauth")
		return model.UserOAuth{}, err
	}

	return user, nil
}

func (u *userOauth[T]) GetByUserID(ctx context.Context, userID string) (model.UserOAuth, error) {
	query := filters.ApplyFilters(u.baseQuery(), filters.IsSelectFilter("user_id", userID))
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building query")
		return model.UserOAuth{}, err
	}

	var user model.UserOAuth
	err = u.Store.GetContext(ctx, &user, q, args...)
	if err != nil {
		log.Err(err).Msg("error getting user oauth")
		return model.UserOAuth{}, err
	}

	return user, nil
}
