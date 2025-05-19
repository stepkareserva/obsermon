package db

import (
	"context"
	"database/sql"
)

// Stmt impl
type sqlStmt struct {
	*sql.Stmt
}

var _ Stmt = (*sqlStmt)(nil)

func (stmt *sqlStmt) ExecContext(ctx context.Context, args ...any) (Result, error) {
	res, err := stmt.Stmt.ExecContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	return sqlResult{Result: res}, nil
}

func (stmt *sqlStmt) QueryContext(ctx context.Context, args ...any) (Rows, error) {
	rows, err := stmt.Stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	return sqlRows{Rows: rows}, rows.Err()
}
