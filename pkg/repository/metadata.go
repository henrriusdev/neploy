package repository

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"neploy.dev/pkg/logger"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/store"
)

type Metadata interface {
	Create(ctx context.Context, metadata model.Metadata) error
	Update(ctx context.Context, metadata model.Metadata) error
	Get(ctx context.Context) (model.Metadata, error)
	GetPrimaryColor(ctx context.Context) (string, error)
	GetSecondaryColor(ctx context.Context) (string, error)
	GetTeamName(ctx context.Context) (string, error)
	GetTeamLogo(ctx context.Context) (string, error)
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
		logger.Error("error building insert query: %v", err)
		return err
	}

	if _, err = m.Store.ExecContext(ctx, query, args...); err != nil {
		logger.Error("error executing insert query: %v", err)
		return err
	}

	return nil
}

func (m *metadata[T]) Update(ctx context.Context, metadata model.Metadata) error {
	q := m.BaseQueryUpdate().Set(metadata).Where(goqu.C("id").Eq(metadata.ID))
	query, args, err := q.ToSQL()
	if err != nil {
		logger.Error("error building update query: %v", err)
		return err
	}

	if _, err = m.Store.ExecContext(ctx, query, args...); err != nil {
		logger.Error("error executing update query: %v", err)
		return err
	}

	return err
}

func (m *metadata[T]) Get(ctx context.Context) (model.Metadata, error) {
	q := m.baseQuery().Limit(1)
	query, args, err := q.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return model.Metadata{}, err
	}

	var metadata model.Metadata
	if err = m.Store.GetContext(ctx, &metadata, query, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return model.Metadata{}, err
	}

	return metadata, nil
}

func (m *metadata[T]) GetPrimaryColor(ctx context.Context) (string, error) {
	q := m.baseQuery().Select("primary_color").Limit(1)
	query, args, err := q.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return "", err
	}

	var primaryColor string
	if err = m.Store.GetContext(ctx, &primaryColor, query, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return "", err
	}

	return primaryColor, nil
}

func (m *metadata[T]) GetSecondaryColor(ctx context.Context) (string, error) {
	q := m.baseQuery().Select("secondary_color").Limit(1)
	query, args, err := q.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return "", err
	}

	var secondaryColor string
	if err = m.Store.GetContext(ctx, &secondaryColor, query, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return "", err
	}

	return secondaryColor, nil
}

func (m *metadata[T]) GetTeamName(ctx context.Context) (string, error) {
	q := m.baseQuery().Select("team_name").Limit(1)
	query, args, err := q.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return "", err
	}

	var teamName string
	if err = m.Store.GetContext(ctx, &teamName, query, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return "", err
	}

	return teamName, nil
}

func (m *metadata[T]) GetTeamLogo(ctx context.Context) (string, error) {
	q := m.baseQuery().Select("logo_url").Limit(1)
	query, args, err := q.ToSQL()
	if err != nil {
		logger.Error("error building select query: %v", err)
		return "", err
	}

	var teamLogo string
	if err = m.Store.GetContext(ctx, &teamLogo, query, args...); err != nil {
		logger.Error("error executing select query: %v", err)
		return "", err
	}

	return teamLogo, nil
}
