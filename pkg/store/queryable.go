package store

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Queryable interface {
	sqlx.Ext
	sqlx.ExecerContext
	sqlx.PreparerContext
	sqlx.QueryerContext
	sqlx.Preparer

	Get(interface{}, string, ...interface{}) error
	GetContext(context.Context, interface{}, string, ...interface{}) error
	Select(interface{}, string, ...interface{}) error
	SelectContext(context.Context, interface{}, string, ...interface{}) error
	MustExec(string, ...interface{}) sql.Result
	MustExecContext(context.Context, string, ...interface{}) sql.Result
	PreparexContext(context.Context, string) (*sqlx.Stmt, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	QueryxContext(context.Context, string, ...interface{}) (*sqlx.Rows, error)
	QueryRowxContext(context.Context, string, ...interface{}) *sqlx.Row
	PrepareNamed(string) (*sqlx.NamedStmt, error)
	PrepareNamedContext(context.Context, string) (*sqlx.NamedStmt, error)
	Preparex(string) (*sqlx.Stmt, error)
	NamedExec(string, interface{}) (sql.Result, error)
	NamedExecContext(context.Context, string, interface{}) (sql.Result, error)
	NamedQuery(string, interface{}) (*sqlx.Rows, error)
	Rebind(string) string
}

type Readable interface {
	Select(interface{}, string, ...interface{}) error
	Get(interface{}, string, ...interface{}) error
}

type Transactable interface {
	MustBegin() Queryable
	Rollback() error
	Commit() error
	SetTx(Queryable)
	Reset(...Transactable)
}
