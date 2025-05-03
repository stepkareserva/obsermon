package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	db *sql.DB
}

func New(dbConn string) (*Database, error) {
	db, err := sql.Open("pgx", dbConn)
	if err != nil {
		return nil, fmt.Errorf("db connection: %w", err)
	}
	return &Database{db: db}, nil
}

func (db *Database) Close() error {
	if db.db == nil {
		return nil
	}
	if err := db.db.Close(); err != nil {
		return fmt.Errorf("database closing: %w", err)
	}
	return nil
}

func (db *Database) Ping() error {
	if db == nil {
		return fmt.Errorf("db not exists")
	}
	if err := db.db.PingContext(context.TODO()); err != nil {
		return fmt.Errorf("db: %w", err)
	}
	return nil
}

func (db *Database) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if db == nil || db.db == nil {
		return nil, fmt.Errorf("database not exists")
	}
	return db.db.ExecContext(ctx, query, args...)
}

func (db *Database) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if db == nil || db.db == nil {
		return nil, fmt.Errorf("database not exists")
	}
	return db.db.QueryContext(ctx, query, args...)
}

func (db *Database) QueryRowContext(ctx context.Context, query string, args ...any) (*sql.Row, error) {
	if db == nil || db.db == nil {
		return nil, fmt.Errorf("database not exists")
	}
	return db.db.QueryRowContext(ctx, query, args...), nil
}
