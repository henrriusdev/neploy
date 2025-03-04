package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository/filters"
	"neploy.dev/pkg/store"
)

type TechStack struct {
	Base[model.TechStack]
}

func NewTechStack(db store.Queryable) *TechStack {
	return &TechStack{Base[model.TechStack]{Store: db, Table: "tech_stacks"}}
}

func (t *TechStack) FindOrCreate(ctx context.Context, name string) (model.TechStack, error) {
	query := t.baseQuery().Where(goqu.Ex{"name": name})
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building find or create query: %v", err)
		return model.TechStack{}, err
	}

	var techStack model.TechStack
	if err := t.Store.GetContext(ctx, &techStack, q, args...); err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Error("error executing find query: %v", err)
		return model.TechStack{}, err
	}

	if techStack.ID != "" {
		return techStack, nil
	}

	insert := t.BaseQueryInsert().Rows(model.TechStack{Name: name}).Returning(goqu.Star())
	q, args, err = insert.ToSQL()
	if err != nil {
		logger.Error("error building insert query: %v", err)
		return model.TechStack{}, err
	}

	if err := t.Store.QueryRowxContext(ctx, q, args...).StructScan(&techStack); err != nil {
		logger.Error("error executing insert query: %v", err)
		return model.TechStack{}, err
	}

	return techStack, nil
}

func (t *TechStack) Insert(ctx context.Context, techStack model.TechStack) error {
	query := t.BaseQueryInsert().Rows(techStack)
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building insert query: %v", err)
		return err
	}

	if _, err := t.Store.ExecContext(ctx, q, args...); err != nil {
		logger.Error("error executing insert query: %v", err)
		return err
	}

	return nil
}

func (t *TechStack) Update(ctx context.Context, id string, techStack model.TechStack) error {
	query := filters.ApplyUpdateFilters(t.BaseQueryUpdate().Set(techStack), filters.IsUpdateFilter("id", id))
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building update query: %v", err)
		return err
	}

	if _, err := t.Store.ExecContext(ctx, q, args...); err != nil {
		logger.Error("error executing update query: %v", err)
		return err
	}

	return nil
}

func (t *TechStack) Delete(ctx context.Context, id string) error {
	query := filters.ApplyUpdateFilters(
		t.BaseQueryUpdate().
			Set(goqu.Record{"deleted_at": goqu.L("CURRENT_TIMESTAMP")}),
		filters.IsUpdateFilter("id", id),
	)

	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building delete query: %v", err)
		return err
	}

	if _, err := t.Store.ExecContext(ctx, q, args...); err != nil {
		logger.Error("error executing delete query: %v", err)
		return err
	}

	return nil
}

func (t *TechStack) GetByID(ctx context.Context, id string) (model.TechStack, error) {
	query := t.baseQuery().Where(goqu.Ex{"id": id})
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building get by id query: %v", err)
		return model.TechStack{}, err
	}

	var techStack model.TechStack
	if err := t.Store.GetContext(ctx, &techStack, q, args...); err != nil {
		logger.Error("error executing get by id query: %v", err)
		return model.TechStack{}, err
	}

	return techStack, nil
}

func (t *TechStack) GetAll(ctx context.Context) ([]model.TechStack, error) {
	query := t.baseQuery().Where(goqu.Ex{"deleted_at": nil})
	q, args, err := query.ToSQL()
	if err != nil {
		logger.Error("error building get all query: %v", err)
		return nil, err
	}

	var techStacks []model.TechStack
	if err := t.Store.SelectContext(ctx, &techStacks, q, args...); err != nil {
		logger.Error("error executing get all query: %v", err)
		return nil, err
	}

	return techStacks, nil
}
