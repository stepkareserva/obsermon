package database

import (
	"context"
	"database/sql"
)

type Database interface {
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
	Queryt(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) (*sql.Row, error)
	ExecTxFn(ctx context.Context, txFn func(ctx context.Context, tx *sql.Tx) error) error

	Ping() error
	Close() error
}
