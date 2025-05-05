package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stepkareserva/obsermon/internal/server/db/handlers"
	"github.com/stepkareserva/obsermon/internal/server/metrics/dbstorage"
)

type Database struct {
	db *sql.DB
}

var _ handlers.Database = (*Database)(nil)
var _ dbstorage.Database = (*Database)(nil)

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

func (db *Database) ExecContext(ctx context.Context, query string, args ...any) (res sql.Result, err error) {
	err = db.sustainedOp(ctx, func(ctx context.Context) error {
		res, err = db.db.ExecContext(ctx, query, args...)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("exec: %w", err)
	}
	return
}

func (db *Database) QueryContext(ctx context.Context, query string, args ...any) (rows *sql.Rows, err error) {
	err = db.sustainedOp(ctx, func(ctx context.Context) error {
		rows, err = db.db.QueryContext(ctx, query, args...)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	return
}

func (db *Database) QueryRowContext(ctx context.Context, query string, args ...any) (row *sql.Row, err error) {
	err = db.sustainedOp(ctx, func(ctx context.Context) error {
		row = db.db.QueryRowContext(ctx, query, args...)
		err = row.Err()
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("query row: %w", err)
	}
	return
}

func (db *Database) ExecTxFn(ctx context.Context, txFn func(ctx context.Context, tx *sql.Tx) error) (err error) {
	err = db.sustainedOp(ctx, func(ctx context.Context) error {
		err = db.execTxFn(ctx, txFn)
		return err
	})
	if err != nil {
		return fmt.Errorf("exec tx: %w", err)
	}
	return
}

func (db *Database) execTxFn(ctx context.Context, txFn func(ctx context.Context, tx *sql.Tx) error) (err error) {
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

func (db *Database) sustainedOp(ctx context.Context, op func(ctx context.Context) error) error {
	if db == nil || db.db == nil {
		return fmt.Errorf("database not exists")
	}

	delays := []time.Duration{
		0,
		1 * time.Second,
		3 * time.Second,
		5 * time.Second,
	}

	for _, delay := range delays {
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			return fmt.Errorf("sustained op timeout")
		case <-timer.C:
			err := op(ctx)
			switch {
			case err == nil:
				return nil
			case !isDBConnectionError(err):
				return fmt.Errorf("sustained op: %w", err)
			}
		}
	}

	return nil
}

func isDBConnectionError(err error) bool {
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}
