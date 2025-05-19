package dbstorage

import (
	"context"
	"errors"
	"fmt"

	"github.com/stepkareserva/obsermon/internal/models"
	"github.com/stepkareserva/obsermon/internal/server/metrics/storage/dbstorage/db"
)

func UpdateGauges(ctx context.Context, uow *UnitOfWork, gauges []models.Gauge) error {
	txFn := func(ctx context.Context, tx db.Tx) (err error) {
		setStmt, err := tx.PrepareContext(ctx, setGaugeQuery)
		if err != nil {
			return fmt.Errorf("prepare set gauge stmt: %w", err)
		}
		defer func() {
			if closeErr := setStmt.Close(); closeErr != nil {
				err = errors.Join(err, fmt.Errorf("close set gauges stmt: %w", closeErr))
			}
		}()

		for _, gauge := range gauges {
			if _, err = setStmt.ExecContext(ctx, gauge.Name, gauge.Value); err != nil {
				return fmt.Errorf("set gauge stmt exec: %w", err)
			}
		}

		return
	}

	return uow.Do(ctx, txFn)
}

type counterStmts struct {
	sel db.Stmt
	upd db.Stmt
	ins db.Stmt
}

func newCounterStmts(ctx context.Context, tx db.Tx) (stmts *counterStmts, err error) {
	closeStmt := func(stmt db.Stmt) {
		if err == nil {
			return
		}
		if closeErr := stmt.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close stmt: %w", closeErr))
		}
	}

	sel, err := tx.PrepareContext(ctx, selectCounterForUpdateQuery)
	if err != nil {
		return nil, fmt.Errorf("prepare sel stmt: %w", err)
	}
	defer func() { closeStmt(sel) }()

	upd, err := tx.PrepareContext(ctx, updateCounterQuery)
	if err != nil {
		return nil, fmt.Errorf("prepare upd stmt: %w", err)
	}
	defer func() { closeStmt(upd) }()

	ins, err := tx.PrepareContext(ctx, insertCounterQuery)
	if err != nil {
		return nil, fmt.Errorf("prepare ins stmt: %w", err)
	}
	defer func() { closeStmt(ins) }()

	return &counterStmts{
		sel: sel,
		upd: upd,
		ins: ins,
	}, nil
}

func (stmts *counterStmts) Close() error {
	var errs []error
	if err := stmts.sel.Close(); err != nil {
		errs = append(errs, fmt.Errorf("close sel stmt: %w", err))
	}
	if err := stmts.upd.Close(); err != nil {
		errs = append(errs, fmt.Errorf("close upd stmt: %w", err))
	}
	if err := stmts.ins.Close(); err != nil {
		errs = append(errs, fmt.Errorf("close ins stmt: %w", err))
	}
	return errors.Join(errs...)
}

func UpdateCounters(ctx context.Context, uow *UnitOfWork, counters []models.Counter) ([]models.Counter, error) {
	var updatedCounters []models.Counter

	txFn := func(ctx context.Context, tx db.Tx) (err error) {
		stmts, err := newCounterStmts(ctx, tx)
		if err != nil {
			return fmt.Errorf("prepare set gauge stmt: %w", err)
		}
		defer func() {
			if closeErr := stmts.Close(); closeErr != nil {
				err = errors.Join(err, fmt.Errorf("close update counters stmts: %w", closeErr))
			}
		}()

		for _, counter := range counters {
			updated, err := updateCounter(ctx, stmts, counter)
			if err != nil {
				return fmt.Errorf("set gauge stmt exec: %w", err)
			}
			updatedCounters = append(updatedCounters, *updated)
		}

		return
	}

	if err := uow.Do(ctx, txFn); err != nil {
		return nil, err
	}

	return updatedCounters, nil
}

func updateCounter(ctx context.Context, stmts *counterStmts, counter models.Counter) (updated *models.Counter, err error) {
	rows, err := stmts.sel.QueryContext(ctx, counter.Name)
	if err != nil {
		return nil, fmt.Errorf("query counters: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("rows closing: %w", err))
		}
	}()

	current, err := ScanCounter(rows)
	if err != nil {
		return nil, fmt.Errorf("scan counters: %w", err)
	}

	if current == nil {
		if _, err = stmts.ins.ExecContext(ctx, counter.Name, counter.Value); err != nil {
			return nil, fmt.Errorf("insert counter: %w", err)
		}
		return &counter, nil
	}

	if err = current.Value.Update(counter.Value); err != nil {
		return nil, fmt.Errorf("update counter value: %w", err)
	}
	if _, err = stmts.upd.ExecContext(ctx, counter.Name); err != nil {
		return nil, fmt.Errorf("update counter: %w", err)
	}
	return current, nil

}
