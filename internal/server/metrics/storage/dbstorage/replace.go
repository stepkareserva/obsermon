package dbstorage

import (
	"context"
	"errors"
	"fmt"

	"github.com/stepkareserva/obsermon/internal/models"
	"github.com/stepkareserva/obsermon/internal/server/metrics/storage/dbstorage/db"
)

func ReplaceGauges(ctx context.Context, uow *UnitOfWork, gauges []models.Gauge) error {
	txFn := func(ctx context.Context, tx db.Tx) (err error) {
		_, err = tx.ExecContext(ctx, clearGaugeQuery)
		if err != nil {
			return fmt.Errorf("clear gauges: %w", err)
		}

		insStmt, err := tx.PrepareContext(ctx, insertGaugeQuery)
		if err != nil {
			return fmt.Errorf("prepare insert gauge stmt: %w", err)
		}
		defer func() {
			if closeErr := insStmt.Close(); closeErr != nil {
				err = errors.Join(err, fmt.Errorf("close ins gauges stmt: %w", closeErr))
			}
		}()

		for _, gauge := range gauges {
			if _, err = insStmt.ExecContext(ctx, gauge.Name, gauge.Value); err != nil {
				return fmt.Errorf("ins gauge stmt exec: %w", err)
			}
		}

		return
	}

	return uow.Do(ctx, txFn)
}

func ReplaceCounters(ctx context.Context, uow *UnitOfWork, counters []models.Counter) error {
	txFn := func(ctx context.Context, tx db.Tx) (err error) {
		_, err = tx.ExecContext(ctx, clearCountersQuery)
		if err != nil {
			return fmt.Errorf("clear counters: %w", err)
		}

		insStmt, err := tx.PrepareContext(ctx, insertCounterQuery)
		if err != nil {
			return fmt.Errorf("prepare insert counter stmt: %w", err)
		}
		defer func() {
			if closeErr := insStmt.Close(); closeErr != nil {
				err = errors.Join(err, fmt.Errorf("close ins counters stmt: %w", closeErr))
			}
		}()

		for _, counter := range counters {
			if _, err = insStmt.ExecContext(ctx, counter.Name, counter.Value); err != nil {
				return fmt.Errorf("ins counter stmt exec: %w", err)
			}
		}

		return
	}

	return uow.Do(ctx, txFn)
}
