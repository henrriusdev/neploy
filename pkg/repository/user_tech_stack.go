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

type UserTechStack struct {
	Base[model.UserTechStack]
}

func NewUserTechStack(db store.Queryable) *UserTechStack {
	return &UserTechStack{Base[model.UserTechStack]{Store: db, Table: "user_tech_stacks"}}
}

func (u *UserTechStack) Insert(ctx context.Context, userTechStack model.UserTechStack) error {
	query := u.BaseQueryInsert().Rows(userTechStack)
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building insert query: %v", err)
		return err
	}

	if _, err := u.Store.ExecContext(ctx, q, args...); err != nil {
		logger.Error("error executing insert query: %v", err)
		return err
	}

	return nil
}

func (u *UserTechStack) Update(ctx context.Context, userTechStack model.UserTechStack) error {
	query := filters.ApplyUpdateFilters(
		u.BaseQueryUpdate().Set(userTechStack),
		filters.IsUpdateFilter("user_id", userTechStack.UserID),
		filters.IsUpdateFilter("tech_stack_id", userTechStack.TechStackID),
	)
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building update query: %v", err)
		return err
	}

	if _, err := u.Store.ExecContext(ctx, q, args...); err != nil {
		logger.Error("error executing update query: %v", err)
		return err
	}

	return nil
}

func (u *UserTechStack) Delete(ctx context.Context, userId, techId string) error {
	query := dialect.Delete(u.Table).Where(goqu.I("user_id").Eq(userId), goqu.I("tech_stack_id").Eq(techId))

	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building delete query: %v", err)
		return err
	}

	if _, err := u.Store.ExecContext(ctx, q, args...); err != nil {
		logger.Error("error executing delete query: %v", err)
		return err
	}

	return nil
}

func (u *UserTechStack) GetByUserID(ctx context.Context, userID string) ([]model.UserTechStack, error) {
	query := filters.ApplyFilters(u.baseQuery(), filters.IsSelectFilter("user_id", userID))
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query for usertechstack: %v", err)
		return nil, err
	}

	var userTechStacks []model.UserTechStack
	if err := u.Store.SelectContext(ctx, &userTechStacks, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return nil, err
	}

	common.AttachSQLToTrace(ctx, q)
	return userTechStacks, nil
}

func (u *UserTechStack) GetAll(ctx context.Context) ([]model.UserTechStack, error) {
	query := u.baseQuery()
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return nil, err
	}

	var userTechStacks []model.UserTechStack
	if err := u.Store.SelectContext(ctx, &userTechStacks, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return nil, err
	}

	common.AttachSQLToTrace(ctx, q)
	return userTechStacks, nil
}
