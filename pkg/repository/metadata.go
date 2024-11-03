package repository

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/store"
)

type Metadata interface {
	Create(ctx context.Context, metadata model.Metadata) error
	Update(ctx context.Context, metadata model.Metadata) error
	Get(ctx context.Context) (model.Metadata, error)
}

type metadata[T any] struct {
	Base[T]
}

func NewMetadata(db store.Queryable) Metadata {
	return &metadata[model.Metadata]{Base[model.Metadata]{Store: db, Table: "metadata"}}
}

func (m *metadata[T]) Create(ctx context.Context, metadata model.Metadata) error {
	q := m.BaseQueryInsert().Rows(metadata)
	query, args, err := q.ToSQL()
	if err != nil {
		return err
	}

	if _, err = m.Store.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return err
}

func (m *metadata[T]) Update(ctx context.Context, metadata model.Metadata) error {
	q := m.BaseQueryUpdate().Set(metadata).Where(goqu.C("id").Eq(metadata.ID))
	query, args, err := q.ToSQL()
	if err != nil {
		return err
	}

	if _, err = m.Store.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return err
}

func (m *metadata[T]) Get(ctx context.Context) (model.Metadata, error) {
	q := m.baseQuery().Limit(1)
	query, args, err := q.ToSQL()
	if err != nil {
		return model.Metadata{}, err
	}

	var metadata model.Metadata
	if err = m.Store.GetContext(ctx, &metadata, query, args...); err != nil {
		return model.Metadata{}, err
	}

	return metadata, nil
}
