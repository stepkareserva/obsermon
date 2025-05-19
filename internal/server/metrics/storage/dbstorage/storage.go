package dbstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/stepkareserva/obsermon/internal/models"
	"github.com/stepkareserva/obsermon/internal/server/metrics/service"
	"github.com/stepkareserva/obsermon/internal/server/metrics/storage/dbstorage/db"
	"go.uber.org/zap"
)

type Storage struct {
	db  db.Db
	uow *UnitOfWork
}

var _ service.Storage = (*Storage)(nil)

func New(dbConn string, log *zap.Logger) (*Storage, error) {
	if log == nil {
		return nil, fmt.Errorf("log not exists")
	}

	db := db.NewSqlDB(dbConn, log)

	retryPolicy := []time.Duration{
		1 * time.Second,
		3 * time.Second,
		5 * time.Second,
	}
	uow := UnitOfWork{db: db, retryPolicy: retryPolicy}

	storage := Storage{
		db:  db,
		uow: &uow,
	}

	//if err := storage.initTables(); err != nil {
	//	return nil, err
	//}

	return &storage, nil
}

func (s *Storage) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("db closing: %w", err)
	}
	return nil
}

func (s *Storage) initTables() error {
	ctx, cancel := context.WithTimeout(context.Background(), SQLOpTimeout)
	defer cancel()

	err := s.uow.Do(ctx, func(ctx context.Context, tx db.Tx) error {
		if _, err := tx.ExecContext(ctx, createCountersQuery); err != nil {
			return fmt.Errorf("counters table creation: %w", err)
		}
		if _, err := tx.ExecContext(ctx, createGaugesQuery); err != nil {
			return fmt.Errorf("gauges table creation: %w", err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("init tables: %w", err)
	}

	return nil
}

func (s *Storage) SetGauge(ctx context.Context, val models.Gauge) error {
	return s.SetGauges(ctx, models.GaugesList{val})
}

func (s *Storage) SetGauges(ctx context.Context, vals models.GaugesList) error {
	if s == nil || s.uow == nil {
		return fmt.Errorf("database not exists")
	}
	return UpdateGauges(ctx, s.uow, vals)
}

func (s *Storage) FindGauge(ctx context.Context, name string) (*models.Gauge, bool, error) {
	if s == nil || s.uow == nil {
		return nil, false, fmt.Errorf("database not exists")
	}
	return SelectGauge(ctx, s.uow, findGaugeQuery, name)
}

func (s *Storage) ListGauges(ctx context.Context) (models.GaugesList, error) {
	if s == nil || s.uow == nil {
		return nil, fmt.Errorf("database not exists")
	}
	return SelectGauges(ctx, s.uow, listGaugesQuery)
}

func (s *Storage) ReplaceGauges(ctx context.Context, val models.GaugesList) error {
	if s == nil || s.uow == nil {
		return fmt.Errorf("database not exists")
	}
	return ReplaceGauges(ctx, s.uow, val)
}

func (s *Storage) UpdateCounter(ctx context.Context, val models.Counter) (*models.Counter, error) {
	updated, err := s.UpdateCounters(ctx, models.CountersList{val})
	if err != nil {
		return nil, err
	}
	return &updated[0], nil
}

func (s *Storage) UpdateCounters(ctx context.Context, vals models.CountersList) (models.CountersList, error) {
	if s == nil || s.uow == nil {
		return nil, fmt.Errorf("database not exists")
	}
	return UpdateCounters(ctx, s.uow, vals)
}

func (s *Storage) FindCounter(ctx context.Context, name string) (*models.Counter, bool, error) {
	if s == nil || s.uow == nil {
		return nil, false, fmt.Errorf("database not exists")
	}
	return SelectCounter(ctx, s.uow, findCounterQuery, name)
}

func (s *Storage) ListCounters(ctx context.Context) (models.CountersList, error) {
	if s == nil || s.uow == nil {
		return nil, fmt.Errorf("database not exists")
	}
	return SelectCounters(ctx, s.uow, listCountersQuery)
}

func (s *Storage) ReplaceCounters(ctx context.Context, val models.CountersList) error {
	if s == nil || s.uow == nil {
		return fmt.Errorf("database not exists")
	}
	return ReplaceCounters(ctx, s.uow, val)
}
