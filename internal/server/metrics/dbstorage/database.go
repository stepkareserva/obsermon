package dbstorage

import (
	"context"
	"database/sql"
)

type Database interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) (*sql.Row, error)

	Begin() (*sql.Tx, error)
}
