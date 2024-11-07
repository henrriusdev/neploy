package repository

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/rs/zerolog/log"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository/filters"
	"neploy.dev/pkg/store"
)

type UserTechStack interface {
	Insert(ctx context.Context, userTechStack model.UserTechStack) error
	Update(ctx context.Context, userTechStack model.UserTechStack) error
	Delete(ctx context.Context, id string) error
	GetByUserID(ctx context.Context, userID string) ([]model.UserTechStack, error)
	GetByTechStackID(ctx context.Context, techStackID string) ([]model.UserTechStack, error)
	GetByUserIDAndTechStackID(ctx context.Context, userID, techStackID string) (model.UserTechStack, error)
	GetAll(ctx context.Context) ([]model.UserTechStack, error)
}

type userTechStack[T any] struct {
	Base[T]
}

func NewUserTechStack(db store.Queryable) UserTechStack {
	return &userTechStack[model.UserTechStack]{Base[model.UserTechStack]{Store: db, Table: "user_tech_stacks"}}
}

func (u *userTechStack[T]) Insert(ctx context.Context, userTechStack model.UserTechStack) error {
	query := u.BaseQueryInsert().Rows(userTechStack)
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building insert query")
		return err
	}

	if _, err := u.Store.ExecContext(ctx, q, args...); err != nil {
		log.Err(err).Msg("error executing insert query")
		return err
	}

	return nil
}

func (u *userTechStack[T]) Update(ctx context.Context, userTechStack model.UserTechStack) error {
	query := filters.ApplyUpdateFilters(
		u.BaseQueryUpdate().Set(userTechStack),
		filters.IsUpdateFilter("user_id", userTechStack.UserID),
		filters.IsUpdateFilter("tech_stack_id", userTechStack.TechStackID),
	)
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building update query")
		return err
	}

	if _, err := u.Store.ExecContext(ctx, q, args...); err != nil {
		log.Err(err).Msg("error executing update query")
		return err
	}

	return nil
}

func (u *userTechStack[T]) Delete(ctx context.Context, id string) error {
	query := filters.ApplyUpdateFilters(
		u.BaseQueryUpdate().
			Set(goqu.Record{"deleted_at": goqu.L("CURRENT_TIMESTAMP")}),
		filters.IsUpdateFilter("id", id),
	)

	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building delete query")
		return err
	}

	if _, err := u.Store.ExecContext(ctx, q, args...); err != nil {
		log.Err(err).Msg("error executing delete query")
		return err
	}

	return nil
}

func (u *userTechStack[T]) GetByUserID(ctx context.Context, userID string) ([]model.UserTechStack, error) {
	query := u.baseQuery().Where(goqu.Ex{"user_id": userID})
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building select query")
		return nil, err
	}

	var userTechStacks []model.UserTechStack
	if err := u.Store.SelectContext(ctx, &userTechStacks, q, args...); err != nil {
		log.Err(err).Msg("error executing select query")
		return nil, err
	}

	return userTechStacks, nil
}

func (u *userTechStack[T]) GetByTechStackID(ctx context.Context, techStackID string) ([]model.UserTechStack, error) {
	query := u.baseQuery().Where(goqu.Ex{"tech_stack_id": techStackID})
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building select query")
		return nil, err
	}

	var userTechStacks []model.UserTechStack
	if err := u.Store.SelectContext(ctx, &userTechStacks, q, args...); err != nil {
		log.Err(err).Msg("error executing select query")
		return nil, err
	}

	return userTechStacks, nil
}

func (u *userTechStack[T]) GetByUserIDAndTechStackID(ctx context.Context, userID, techStackID string) (model.UserTechStack, error) {
	query := u.baseQuery().Where(goqu.Ex{"user_id": userID, "tech_stack_id": techStackID})
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building select query")
		return model.UserTechStack{}, err
	}

	var userTechStack model.UserTechStack
	if err := u.Store.GetContext(ctx, &userTechStack, q, args...); err != nil {
		log.Err(err).Msg("error executing select query")
		return model.UserTechStack{}, err
	}

	return userTechStack, nil
}

func (u *userTechStack[T]) GetAll(ctx context.Context) ([]model.UserTechStack, error) {
	query := u.baseQuery()
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building select query")
		return nil, err
	}

	var userTechStacks []model.UserTechStack
	if err := u.Store.SelectContext(ctx, &userTechStacks, q, args...); err != nil {
		log.Err(err).Msg("error executing select query")
		return nil, err
	}

	return userTechStacks, nil
}
