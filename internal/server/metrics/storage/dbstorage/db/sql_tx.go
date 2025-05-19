package db

import (
	"context"
	"database/sql"
)

// Tx impl
type sqlTx struct {
	*sql.Tx
}

var _ Tx = (*sqlTx)(nil)

func (tx *sqlTx) PrepareContext(ctx context.Context, query string) (Stmt, error) {
	stmt, err := tx.Tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return &sqlStmt{Stmt: stmt}, nil
}

func (tx *sqlTx) ExecContext(ctx context.Context, query string, args ...any) (Result, error) {
	res, err := tx.Tx.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return sqlResult{Result: res}, nil
}

func (tx *sqlTx) QueryContext(ctx context.Context, query string, args ...any) (Rows, error) {
	rows, err := tx.Tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return sqlRows{Rows: rows}, nil
}

func (tx *sqlTx) Commit() error {
	return tx.Tx.Commit()
}

func (tx *sqlTx) Rollback() error {
	return tx.Tx.Rollback()
}
