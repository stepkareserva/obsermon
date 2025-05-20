package dbstorage

import (
	"context"
	"errors"
	"fmt"

	"github.com/stepkareserva/obsermon/internal/models"
	"github.com/stepkareserva/obsermon/internal/server/metrics/storage/dbstorage/db"
)

func SelectCounters(ctx context.Context, uow *UnitOfWork, query string, args ...any) (models.CountersList, error) {
	var counters models.CountersList

	txFn := func(ctx context.Context, tx db.Tx) (err error) {
		rows, err := tx.QueryContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("query counters: %w", err)
		}
		defer func() {
			if closeErr := rows.Close(); closeErr != nil {
				err = errors.Join(err, fmt.Errorf("rows closing: %w", err))
			}
		}()

		counters, err = ScanCounters(rows)
		if err != nil {
			return fmt.Errorf("scan counters: %w", err)
		}
		return nil
	}

	if err := uow.Do(ctx, txFn); err != nil {
		return nil, err
	}
	return counters, nil
}

func SelectCounter(ctx context.Context, uow *UnitOfWork, query string, args ...any) (*models.Counter, bool, error) {
	counters, err := SelectCounters(ctx, uow, query, args...)
	if err != nil {
		return nil, false, fmt.Errorf("select counter: %w", err)
	}
	switch len(counters) {
	case 0:
		return nil, false, nil
	case 1:
		counter := counters[0]
		return &counter, true, nil
	default:
		return nil, false, fmt.Errorf("more than one counters with the same name")
	}
}

func SelectGauges(ctx context.Context, uow *UnitOfWork, query string, args ...any) (models.GaugesList, error) {
	var gauges models.GaugesList

	txFn := func(ctx context.Context, tx db.Tx) (err error) {
		rows, err := tx.QueryContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("query gauges: %w", err)
		}
		defer func() {
			if closeErr := rows.Close(); closeErr != nil {
				err = errors.Join(err, fmt.Errorf("rows closing: %w", err))
			}
		}()

		gauges, err = ScanGauges(rows)
		if err != nil {
			return fmt.Errorf("scan gauges: %w", err)
		}
		return nil
	}

	if err := uow.Do(ctx, txFn); err != nil {
		return nil, err
	}
	return gauges, nil
}

func SelectGauge(ctx context.Context, uow *UnitOfWork, query string, args ...any) (*models.Gauge, bool, error) {
	gauges, err := SelectGauges(ctx, uow, query, args...)
	if err != nil {
		return nil, false, fmt.Errorf("select gauge: %w", err)
	}
	switch len(gauges) {
	case 0:
		return nil, false, nil
	case 1:
		gauge := gauges[0]
		return &gauge, true, nil
	default:
		return nil, false, fmt.Errorf("more than one gauges with the same name")
	}
}
