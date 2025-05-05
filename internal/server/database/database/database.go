package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stepkareserva/obsermon/internal/server/interfaces/database"
)

type Database struct {
	db *sql.DB
}

var _ database.Database = (*Database)(nil)

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

func (db *Database) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if db == nil || db.db == nil {
		return nil, fmt.Errorf("db not exists")
	}
	res, err := db.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec: %w", err)
	}
	return res, nil
}

func (db *Database) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if db == nil || db.db == nil {
		return nil, fmt.Errorf("db not exists")
	}
	rows, err := db.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	return rows, nil
}

func (db *Database) QueryRow(ctx context.Context, query string, args ...any) (*sql.Row, error) {
	if db == nil || db.db == nil {
		return nil, fmt.Errorf("db not exists")
	}
	row := db.db.QueryRowContext(ctx, query, args...)
	return row, row.Err()
}

func (db *Database) ExecTxFn(ctx context.Context, txFn database.TxFn) (err error) {
	if db == nil || db.db == nil {
		return fmt.Errorf("database not exists")
	}

	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = rbErr
			}
		} else {
			err = tx.Commit()
		}
	}()

	err = txFn(ctx, tx)
	if err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
