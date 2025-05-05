package database

import (
	"context"
	"database/sql"
)

type TxFn = func(ctx context.Context, tx *sql.Tx) error

type Database interface {
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
	Query(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) (*sql.Row, error)
	ExecTxFn(ctx context.Context, txFn TxFn) error

	Ping() error
	Close() error
}
