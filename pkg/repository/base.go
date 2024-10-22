package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"neploy.dev/pkg/repository/filters"
	"neploy.dev/pkg/store"
)

type Repositories struct {
	User User
}

var (
	ErrNotFound = errors.New("not found")
	dialect     = goqu.Dialect("postgres")
)

type Base[T any] struct {
	Store store.Queryable
	DB    store.Queryable
	Table string
}

// Transaction helpers for managing transactions
func (b *Base[T]) MustBegin() store.Queryable {
	db := b.Store.(*sqlx.DB)
	b.DB = db
	t := db.MustBegin()
	b.Store = t
	return t
}

func (b *Base[T]) Rollback() {
	t := b.Store.(*sqlx.Tx)
	t.Rollback()
	b.Reset()
}

func (b *Base[T]) Commit() error {
	t := b.Store.(*sqlx.Tx)
	err := t.Commit()
	if err != nil {
		return err
	}
	return err
}

func (b *Base[T]) SetTx(t store.Queryable) {
	b.DB = b.Store
	b.Store = t
}

func (b *Base[T]) Reset(repos ...store.Transactable) {
	b.Store = b.DB
	for _, v := range repos {
		v.Reset()
	}
}

// Helper function to generate a query builder with base filters (soft delete)
func (b *Base[T]) baseQuery(aliases ...string) *goqu.SelectDataset {
	alias := ""
	if len(aliases) > 0 {
		alias = aliases[0] // Use the first one
	}
	var tableExp exp.Expression
	table := goqu.T(b.Table)

	if alias != "" {
		tableExp = table.As(alias) // Apply alias if present
	} else {
		tableExp = table // Use table directly without alias
	}

	softDeleteColumn := "deleted_at"
	if alias != "" {
		softDeleteColumn = alias + ".deleted_at"
	}
	return dialect.From(tableExp).Where(goqu.I(softDeleteColumn).IsNull())
}

func (b *Base[T]) BaseQueryUpdate() *goqu.UpdateDataset {
	return dialect.Update(b.Table)
}

func (b *Base[T]) BaseQueryInsert() *goqu.InsertDataset {
	return dialect.Insert(b.Table)
}

// List retrieves a list of records with pagination and soft delete filtering
func (b Base[T]) List(ctx context.Context, limit int, offset int) (list []T, err error) {
	if limit == 0 {
		limit = 1000 // Default limit
	}

	query := b.baseQuery()
	query = filters.ApplyFilters(query,
		filters.LimitOffsetFilter(uint(limit), uint(offset)), // Pagination filter
	)

	sq, args, err := query.ToSQL()
	if err != nil {
		return list, err
	}

	err = b.Store.SelectContext(ctx, &list, sq, args...)
	if err == sql.ErrNoRows {
		return list, nil
	}
	return list, err
}

// GetByID retrieves a record by its ID and applies soft delete filtering
func (b Base[T]) GetByID(ctx context.Context, id string) (m T, err error) {
	query := b.baseQuery()
	query = filters.ApplyFilters(query,
		filters.GenericColumnSelectFilter("id", id, ""), // Filter by ID
		filters.LimitOffsetFilter(1, 0),                 // Limit 1
	)

	sq, args, err := query.ToSQL()
	if err != nil {
		return m, err
	}

	err = b.Store.GetContext(ctx, &m, sq, args...)
	if err == sql.ErrNoRows {
		return m, ErrNotFound
	}
	return m, err
}

// GetByUserID retrieves a record by user ID with soft delete filtering
func (b Base[T]) GetByUserID(ctx context.Context, userID string) (m T, err error) {
	query := b.baseQuery()
	query = filters.ApplyFilters(query,
		filters.GenericColumnSelectFilter("user_id", userID, ""), // Filter by user_id
		filters.LimitOffsetFilter(1, 0),                          // Limit 1
	)

	sq, args, err := query.ToSQL()
	if err != nil {
		return m, err
	}

	err = b.Store.GetContext(ctx, &m, sq, args...)
	if err == sql.ErrNoRows {
		return m, ErrNotFound
	}
	return m, err
}

// ListByUserID retrieves a paginated list of records for a specific user ID
func (b Base[T]) ListByUserID(ctx context.Context, userID string, limit int, offset int) ([]T, error) {
	if limit == 0 {
		limit = 1000
	}

	query := b.baseQuery()
	query = filters.ApplyFilters(query,
		filters.GenericColumnSelectFilter("user_id", userID, ""), // Filter by user ID
		filters.LimitOffsetFilter(uint(limit), uint(offset)),     // Pagination
	)

	sq, args, err := query.ToSQL()
	if err != nil {
		return nil, err
	}

	var list []T
	err = b.Store.SelectContext(ctx, &list, sq, args...)
	if err == sql.ErrNoRows {
		return list, nil
	}
	return list, err
}

// SoftDelete marks a record as deleted by setting the deleted_at timestamp
func (b Base[T]) SoftDelete(ctx context.Context, id string) error {
	now := time.Now()

	query := goqu.Update(b.Table).
		Set(goqu.Record{"deleted_at": now}).
		Where(goqu.C("id").Eq(id), goqu.C("deleted_at").IsNull())

	sql, args, err := query.ToSQL()
	if err != nil {
		return err
	}

	_, err = b.Store.ExecContext(ctx, sql, args...)
	return err
}

// IsDeleted checks if a record is soft-deleted
func (b Base[T]) IsDeleted(ctx context.Context, id string) (bool, error) {
	query := goqu.From(b.Table).
		Select(goqu.COUNT("*")).
		Where(goqu.C("id").Eq(id), goqu.C("deleted_at").IsNotNull())

	sql, args, err := query.ToSQL()
	if err != nil {
		return false, err
	}

	var count int
	err = b.Store.GetContext(ctx, &count, sql, args...)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Rebind is used to rebind query placeholders (for databases like Postgres)
func (b *Base[T]) Rebind(query string) string {
	db, ok := b.Store.(*sqlx.DB)
	if !ok {
		tx, ok := b.Store.(*sqlx.Tx)
		if !ok {
			log.Error().Msg("Store is not a *sqlx.DB or *sqlx.Tx")
			return ""
		}
		return tx.Rebind(query)
	}
	return db.Rebind(query)
}
