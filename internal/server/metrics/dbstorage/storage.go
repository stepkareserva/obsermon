package dbstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/stepkareserva/obsermon/internal/models"
	"github.com/stepkareserva/obsermon/internal/server/metrics/service"
	"go.uber.org/zap"
)

type Storage struct {
	db  Database
	log *zap.Logger
}

var _ service.Storage = (*Storage)(nil)

const (
	CountersTable = "counters"
	GaugesTable   = "gauges"

	NameColumn  = "name"
	ValueColumn = "value"

	SqlOpTimeout = 5 * time.Second
)

var (
	queryReplacer = strings.NewReplacer(
		"{counters}", CountersTable,
		"{gauges}", GaugesTable,
		"{name}", NameColumn,
		"{value}", ValueColumn)

	createCountersQuery = queryReplacer.Replace(`
		CREATE TABLE IF NOT EXISTS {counters} (
			{name} TEXT PRIMARY KEY,
			{value} BIGINT NOT NULL
		)`)

	createGaugesQuery = queryReplacer.Replace(`
		CREATE TABLE IF NOT EXISTS {gauges} (
			{name} TEXT PRIMARY KEY,
			{value} DOUBLE PRECISION NOT NULL
		)`)

	setCounterQuery = queryReplacer.Replace(`
		INSERT
			INTO {counters} ({name}, {value})
			VALUES ($1, $2)
		ON CONFLICT ({name})
			DO UPDATE SET {value} = EXCLUDED.{value};
		`)

	findCounterQuery = queryReplacer.Replace(`
		SELECT {value}
			FROM {counters}
			WHERE {name} = $1;
		`)

	listCountersQuery = queryReplacer.Replace(`
		SELECT {name}, {value}
			FROM {counters}
		`)

	setGaugeQuery = queryReplacer.Replace(`
		INSERT
			INTO {gauges} ({name}, {value})
			VALUES ($1, $2)
		ON CONFLICT ({name})
			DO UPDATE SET {value} = EXCLUDED.{value};
		`)

	findGaugeQuery = queryReplacer.Replace(`
		SELECT {value}
			FROM {gauges}
			WHERE {name} = $1;
		`)

	listGaugesQuery = queryReplacer.Replace(`
		SELECT {name}, {value}
			FROM {gauges}
		`)
)

func New(db Database, log *zap.Logger) (*Storage, error) {
	if db == nil {
		return nil, fmt.Errorf("database not exists")
	}
	if log == nil {
		log = zap.NewNop()
	}
	storage := Storage{
		db:  db,
		log: log,
	}

	if err := storage.initTables(); err != nil {
		return nil, err
	}

	return &storage, nil
}

func (s *Storage) initTables() error {
	// create counters table
	createCountersCtx, _ := context.WithTimeout(context.Background(), SqlOpTimeout)
	if _, err := s.db.ExecContext(createCountersCtx, createCountersQuery); err != nil {
		return fmt.Errorf("counters table creation: %w", err)
	}

	// create gauges table
	createGaugesCtx, _ := context.WithTimeout(context.Background(), SqlOpTimeout)
	if _, err := s.db.ExecContext(createGaugesCtx, createGaugesQuery); err != nil {
		return fmt.Errorf("gauges table creation: %w", err)
	}

	return nil
}

func (s *Storage) SetGauge(val models.Gauge) error {
	if err := s.setMetric(setGaugeQuery, val.Name, val.Value); err != nil {
		return fmt.Errorf("set gauge: %w", err)
	}
	return nil
}

func (s *Storage) FindGauge(name string) (*models.Gauge, bool, error) {
	var value models.GaugeValue
	exists, err := s.findMetric(findGaugeQuery, name, &value)
	if err != nil {
		return nil, false, fmt.Errorf("find gauge: %w", err)
	}
	if !exists {
		return nil, false, nil
	}
	return &models.Gauge{
		Name:  name,
		Value: value,
	}, true, nil
}

func (s *Storage) ListGauges() (models.GaugesList, error) {
	rows, err := s.queryRows(listGaugesQuery)
	if err != nil {
		return nil, fmt.Errorf("query gauges list: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			s.log.Error("gauges rows closing", zap.Error(err))
		}
	}()

	var gauges models.GaugesList
	for rows.Next() {
		var gauge models.Gauge
		if err := rows.Scan(&gauge.Name, &gauge.Value); err != nil {
			return nil, err
		}
		gauges = append(gauges, gauge)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("gauges rows reading: %w", err)
	}

	return gauges, nil
}

func (s *Storage) ReplaceGauges(val models.GaugesList) error {
	return nil
}

func (s *Storage) SetCounter(val models.Counter) error {
	if err := s.setMetric(setCounterQuery, val.Name, val.Value); err != nil {
		return fmt.Errorf("set counter: %w", err)
	}
	return nil
}

func (s *Storage) FindCounter(name string) (*models.Counter, bool, error) {
	var value models.CounterValue
	exists, err := s.findMetric(findCounterQuery, name, &value)
	if err != nil {
		return nil, false, fmt.Errorf("find gauge: %w", err)
	}
	if !exists {
		return nil, false, nil
	}
	return &models.Counter{
		Name:  name,
		Value: value,
	}, true, nil
}

func (s *Storage) ListCounters() (models.CountersList, error) {
	rows, err := s.queryRows(listCountersQuery)
	if err != nil {
		return nil, fmt.Errorf("query counters list: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			s.log.Error("counters rows closing", zap.Error(err))
		}
	}()

	var counters models.CountersList
	for rows.Next() {
		var counter models.Counter
		if err := rows.Scan(&counter.Name, &counter.Value); err != nil {
			return nil, err
		}
		counters = append(counters, counter)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("counters rows reading: %w", err)
	}

	return counters, nil
}

func (s *Storage) ReplaceCounters(val models.CountersList) error {
	return nil
}

func (s *Storage) setMetric(query, name string, value any) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("database not exists")
	}

	ctx, _ := context.WithTimeout(context.Background(), SqlOpTimeout)
	if _, err := s.db.ExecContext(ctx, query, name, value); err != nil {
		return fmt.Errorf("set metric: %w", err)
	}

	return nil
}

func (s *Storage) findMetric(query, name string, value any) (bool, error) {
	if s == nil || s.db == nil {
		return false, fmt.Errorf("database not exists")
	}

	ctx, _ := context.WithTimeout(context.Background(), SqlOpTimeout)
	row, err := s.db.QueryRowContext(ctx, query, name)
	if err != nil {
		return false, fmt.Errorf("find metric: %w", err)
	}

	if err := row.Scan(value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		} else {
			return false, fmt.Errorf("metric parsing: %w", err)
		}
	}

	return true, nil
}

func (s *Storage) queryRows(query string) (*sql.Rows, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("database not exists")
	}
	ctx, _ := context.WithTimeout(context.Background(), SqlOpTimeout)
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query rows: %w", err)
	}
	return rows, nil
}
