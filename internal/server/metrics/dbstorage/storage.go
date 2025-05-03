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

type QueryTemplates struct {
	setCounter   string
	findCounter  string
	listCounters string

	setGauge   string
	findGauge  string
	listGauges string
}

type Storage struct {
	db      Database
	queries QueryTemplates
	log     *zap.Logger
}

var _ service.Storage = (*Storage)(nil)

const (
	CountersTable = "counters"
	GaugesTable   = "gauges"

	NameColumn  = "name"
	ValueColumn = "value"

	SqlOpTimeout = 5 * time.Second
)

func New(db Database, log *zap.Logger) (*Storage, error) {
	if db == nil {
		return nil, fmt.Errorf("database not exists")
	}
	if log == nil {
		log = zap.NewNop()
	}
	storage := Storage{
		db:      db,
		log:     log,
		queries: getQueryTemplates(),
	}

	if err := storage.initTables(); err != nil {
		return nil, err
	}

	return &storage, nil
}

func getQueryTemplates() QueryTemplates {
	var queries QueryTemplates

	replacer := strings.NewReplacer(
		"{counters}", CountersTable,
		"{gauges}", GaugesTable,
		"{name}", NameColumn,
		"{value}", ValueColumn)

	queries.setCounter = replacer.Replace(`
	INSERT
		INTO {counters} ({name}, {value})
		VALUES ($1, $2)
	ON CONFLICT ({name})
		DO UPDATE SET {value} = EXCLUDED.{value};
	`)

	queries.findCounter = replacer.Replace(`
	SELECT {value}
		FROM {counters}
		WHERE {name} = $1;
	`)

	queries.listCounters = replacer.Replace(`
	SELECT {name}, {value}
		FROM {counters}
	`)

	queries.setGauge = replacer.Replace(`
	INSERT
		INTO {gauges} ({name}, {value})
		VALUES ($1, $2)
	ON CONFLICT ({name})
		DO UPDATE SET {value} = EXCLUDED.{value};
	`)

	queries.findGauge = replacer.Replace(`
	SELECT {value}
		FROM {gauges}
		WHERE {name} = $1;
	`)

	queries.listGauges = replacer.Replace(`
	SELECT {name}, {value}
		FROM {gauges}
	`)

	return queries
}

func (s *Storage) initTables() error {
	// create counters table
	createCountersQuery := fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s (
        %s TEXT PRIMARY KEY,
        %s BIGINT NOT NULL
    );
	`, CountersTable, NameColumn, ValueColumn)
	createCountersCtx, _ := context.WithTimeout(context.Background(), SqlOpTimeout)
	if _, err := s.db.ExecContext(createCountersCtx, createCountersQuery); err != nil {
		return fmt.Errorf("counters table creation: %w", err)
	}

	// create gauges table
	createGaugesQuery := fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s (
        %s TEXT PRIMARY KEY,
        %s DOUBLE PRECISION NOT NULL
    );
	`,
		GaugesTable,
		NameColumn,
		ValueColumn,
	)

	createGaugesCtx, _ := context.WithTimeout(context.Background(), SqlOpTimeout)
	if _, err := s.db.ExecContext(createGaugesCtx, createGaugesQuery); err != nil {
		return fmt.Errorf("gauges table creation: %w", err)
	}

	return nil
}

func (s *Storage) SetGauge(val models.Gauge) error {
	if err := s.setMetric(s.queries.setGauge, val.Name, val.Value); err != nil {
		return fmt.Errorf("set gauge: %w", err)
	}
	return nil
}

func (s *Storage) FindGauge(name string) (*models.Gauge, bool, error) {
	var value models.GaugeValue
	exists, err := s.findMetric(s.queries.findGauge, name, &value)
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
	return nil, nil
}

func (s *Storage) ReplaceGauges(val models.GaugesList) error {
	return nil
}

func (s *Storage) SetCounter(val models.Counter) error {
	if err := s.setMetric(s.queries.setCounter, val.Name, val.Value); err != nil {
		return fmt.Errorf("set counter: %w", err)
	}
	return nil
}

func (s *Storage) FindCounter(name string) (*models.Counter, bool, error) {
	var value models.CounterValue
	exists, err := s.findMetric(s.queries.findCounter, name, &value)
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
	return nil, nil
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
