package db

import (
	"context"
)

//go:generate mockgen -source=$GOFILE -destination=../../../../mocks/mock_database.go -package=mocks

type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
	Close() error
}

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

type Stmt interface {
	ExecContext(ctx context.Context, args ...any) (Result, error)
	QueryContext(ctx context.Context, args ...any) (Rows, error)
	Close() error
}

type Tx interface {
	PrepareContext(ctx context.Context, query string) (Stmt, error)
	ExecContext(ctx context.Context, query string, args ...any) (Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (Rows, error)
	Commit() error
	Rollback() error
}

type Db interface {
	BeginTx(ctx context.Context) (Tx, error)
	PingContext(ctx context.Context) error
	Close() error
}
