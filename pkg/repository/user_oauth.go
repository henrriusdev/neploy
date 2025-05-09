package repository

import (
	"context"
	"neploy.dev/pkg/common"

	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository/filters"
	"neploy.dev/pkg/store"
)

type UserOAuth struct {
	Base[model.UserOAuth]
}

func NewUserOauth(db store.Queryable) *UserOAuth {
	return &UserOAuth{Base[model.UserOAuth]{Store: db, Table: "user_oauth"}}
}

func (u *UserOAuth) Insert(ctx context.Context, user model.UserOAuth) error {
	query := u.BaseQueryInsert().Rows(user)
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building query: %v", err)
		return err
	}

	_, err = u.Store.ExecContext(ctx, q, args...)
	if err != nil {
		logger.Error("error inserting user oauth: %v", err)
		return err
	}

	common.AttachSQLToTrace(ctx, q)
	return nil
}

func (u *UserOAuth) GetByOAuthID(ctx context.Context, oauthID string) (model.UserOAuth, error) {
	query := filters.ApplyFilters(u.baseQuery(), filters.IsSelectFilter("oauth_id", oauthID))
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building query: %v", err)
		return model.UserOAuth{}, err
	}

	var user model.UserOAuth
	err = u.Store.GetContext(ctx, &user, q, args...)
	if err != nil {
		logger.Error("error getting user oauth: %v", err)
		return model.UserOAuth{}, err
	}

	return user, nil
}

func (u *UserOAuth) GetByUserID(ctx context.Context, userID string) (model.UserOAuth, error) {
	query := filters.ApplyFilters(u.baseQuery(), filters.IsSelectFilter("user_id", userID))
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building query: %v", err)
		return model.UserOAuth{}, err
	}

	var user model.UserOAuth
	err = u.Store.GetContext(ctx, &user, q, args...)
	if err != nil {
		logger.Error("error getting user oauth: %v", err)
		return model.UserOAuth{}, err
	}

	common.AttachSQLToTrace(ctx, q)
	return user, nil
}
