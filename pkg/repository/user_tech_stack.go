package repository

import (
	"context"

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

	return userTechStacks, nil
}

func (u *UserTechStack) GetByTechStackID(ctx context.Context, techStackID string) ([]model.UserTechStack, error) {
	query := u.baseQuery("ut").
		Select(
			goqu.I("ut.*"),
			goqu.L(`"u"."id" AS "user.id"`),
			goqu.L(`"u"."first_name" AS "user.first_name"`),
			goqu.L(`"u"."last_name" AS "user.last_name"`),
			goqu.L(`"u"."email" AS "user.email"`),
			goqu.L(`"t"."id" AS "tech_stack.id"`),
			goqu.L(`"t"."name" AS "tech_stack.name"`),
			goqu.L(`"t"."description" AS "tech_stack.description"`),
		).
		LeftJoin(
			goqu.T("users").As("u"),
			goqu.On(goqu.I("u.id").Eq(goqu.I("ut.user_id"))),
		).
		LeftJoin(
			goqu.T("tech_stacks").As("t"),
			goqu.On(goqu.I("t.id").Eq(goqu.I("ut.tech_stack_id"))),
		).
		Where(goqu.I("t.id").Eq(techStackID))
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

	return userTechStacks, nil
}

func (u *UserTechStack) GetByUserIDAndTechStackID(ctx context.Context, userID, techStackID string) (model.UserTechStack, error) {
	query := u.baseQuery().Where(goqu.Ex{"user_id": userID, "tech_stack_id": techStackID})
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return model.UserTechStack{}, err
	}

	var userTechStack model.UserTechStack
	if err := u.Store.GetContext(ctx, &userTechStack, q, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return model.UserTechStack{}, err
	}

	return userTechStack, nil
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

	return userTechStacks, nil
}
