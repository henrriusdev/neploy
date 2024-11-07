package repository

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/rs/zerolog/log"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository/filters"
	"neploy.dev/pkg/store"
)

type TechStack interface {
	Insert(ctx context.Context, techStack model.TechStack) error
	Update(ctx context.Context, techStack model.TechStack) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (model.TechStack, error)
	GetAll(ctx context.Context) ([]model.TechStack, error)
}

type techStack[T any] struct {
	Base[T]
}

func NewTechStack(db store.Queryable) TechStack {
	return &techStack[model.TechStack]{Base[model.TechStack]{Store: db, Table: "tech_stacks"}}
}

func (t *techStack[T]) Insert(ctx context.Context, techStack model.TechStack) error {
	query := t.BaseQueryInsert().Rows(techStack)
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building insert query")
		return err
	}

	if _, err := t.Store.ExecContext(ctx, q, args...); err != nil {
		log.Err(err).Msg("error executing insert query")
		return err
	}

	return nil
}

func (t *techStack[T]) Update(ctx context.Context, techStack model.TechStack) error {
	query := filters.ApplyUpdateFilters(t.BaseQueryUpdate().Set(techStack), filters.IsUpdateFilter("id", techStack.ID))
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building update query")
		return err
	}

	if _, err := t.Store.ExecContext(ctx, q, args...); err != nil {
		log.Err(err).Msg("error executing update query")
		return err
	}

	return nil
}

func (t *techStack[T]) Delete(ctx context.Context, id string) error {
	query := filters.ApplyUpdateFilters(
		t.BaseQueryUpdate().
			Set(goqu.Record{"deleted_at": goqu.L("CURRENT_TIMESTAMP")}),
		filters.IsUpdateFilter("id", id),
	)

	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building delete query")
		return err
	}

	if _, err := t.Store.ExecContext(ctx, q, args...); err != nil {
		log.Err(err).Msg("error executing delete query")
		return err
	}

	return nil
}

func (t *techStack[T]) GetByID(ctx context.Context, id string) (model.TechStack, error) {
	query := t.baseQuery().Where(goqu.Ex{"id": id})
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building get by id query")
		return model.TechStack{}, err
	}

	var techStack model.TechStack
	if err := t.Store.GetContext(ctx, &techStack, q, args...); err != nil {
		log.Err(err).Msg("error executing get by id query")
		return model.TechStack{}, err
	}

	return techStack, nil
}

func (t *techStack[T]) GetAll(ctx context.Context) ([]model.TechStack, error) {
	query := t.baseQuery().Where(goqu.Ex{"deleted_at": nil})
	q, args, err := query.ToSQL()
	if err != nil {
		log.Err(err).Msg("error building get all query")
		return nil, err
	}

	var techStacks []model.TechStack
	if err := t.Store.SelectContext(ctx, &techStacks, q, args...); err != nil {
		log.Err(err).Msg("error executing get all query")
		return nil, err
	}

	return techStacks, nil
}
