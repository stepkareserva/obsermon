package sustained

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/stepkareserva/obsermon/internal/server/interfaces/database"
)

type Database struct {
	database.Database
}

func New(db database.Database) (*Database, error) {
	if db == nil {
		return nil, fmt.Errorf("database not exists")
	}
	return &Database{Database: db}, nil
}

func (db *Database) Exec(ctx context.Context, query string, args ...any) (res sql.Result, err error) {
	err = db.sustainedOp(ctx, func(ctx context.Context) (opErr error) {
		res, opErr = db.Database.Exec(ctx, query, args...)
		return
	})
	if err != nil {
		return nil, fmt.Errorf("db exec: %w", err)
	}
	return
}

func (db *Database) Query(ctx context.Context, query string, args ...any) (rows *sql.Rows, err error) {
	err = db.sustainedOp(ctx, func(ctx context.Context) (opErr error) {
		rows, opErr = db.Database.Query(ctx, query, args...)
		// for static code analyser, idk why, may be false alarm
		if rows != nil && rows.Err() != nil {
			opErr = rows.Err()
		}
		return
	})
	if err != nil {
		return nil, fmt.Errorf("db query rows: %w", err)
	}
	return
}

func (db *Database) QueryRow(ctx context.Context, query string, args ...any) (row *sql.Row, err error) {
	err = db.sustainedOp(ctx, func(ctx context.Context) (opErr error) {
		row, opErr = db.Database.QueryRow(ctx, query, args...)
		return
	})
	if err != nil {
		return nil, fmt.Errorf("db query row: %w", err)
	}
	return
}

func (db *Database) ExecTxFn(ctx context.Context, txFn database.TxFn) (err error) {
	err = db.sustainedOp(ctx, func(ctx context.Context) (opErr error) {
		opErr = db.Database.ExecTxFn(ctx, txFn)
		return
	})
	if err != nil {
		return fmt.Errorf("db exec transaction: %w", err)
	}
	return
}

func (db *Database) sustainedOp(ctx context.Context, op func(ctx context.Context) error) error {
	if db == nil || db.Database == nil {
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
			case !db.couldRetryOperation(err):
				return fmt.Errorf("db sustained op: %w", err)
			}
		}
	}

	return nil
}

func (db *Database) couldRetryOperation(err error) bool {
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}
