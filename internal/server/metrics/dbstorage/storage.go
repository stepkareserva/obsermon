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

	clearCountersQuery = queryReplacer.Replace(`
		DELETE FROM {counters}
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

	clearGaugeQuery = queryReplacer.Replace(`
		DELETE FROM {gauges}
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
	if err := setMetric(s, setGaugeQuery, val.Name, val.Value); err != nil {
		return fmt.Errorf("set gauge: %w", err)
	}
	return nil
}

func (s *Storage) FindGauge(name string) (*models.Gauge, bool, error) {
	value, exists, err := findMetric[models.GaugeValue](s, findGaugeQuery, name)
	if err != nil {
		return nil, false, fmt.Errorf("find gauge: %w", err)
	}
	if !exists {
		return nil, false, nil
	}
	return &models.Gauge{
		Name:  name,
		Value: *value,
	}, true, nil
}

func (s *Storage) ListGauges() (models.GaugesList, error) {
	names, values, err := listMetrics[models.GaugeValue](s, listGaugesQuery)
	if err != nil {
		return nil, fmt.Errorf("list gauges: %w", err)
	}
	gauges, err := zipGauges(names, values)
	if err != nil {
		return nil, fmt.Errorf("zipping gauges: %w", err)
	}
	return gauges, nil
}

func (s *Storage) ReplaceGauges(val models.GaugesList) error {
	names, values := unzipGauges(val)
	if err := replaceMetrics(s, clearGaugeQuery, setGaugeQuery, names, values); err != nil {
		return fmt.Errorf("replace gauges: %w", err)
	}
	return nil
}

func (s *Storage) SetCounter(val models.Counter) error {
	if err := setMetric(s, setCounterQuery, val.Name, val.Value); err != nil {
		return fmt.Errorf("set counter: %w", err)
	}
	return nil
}

func (s *Storage) FindCounter(name string) (*models.Counter, bool, error) {
	value, exists, err := findMetric[models.CounterValue](s, findCounterQuery, name)
	if err != nil {
		return nil, false, fmt.Errorf("find gauge: %w", err)
	}
	if !exists {
		return nil, false, nil
	}
	return &models.Counter{
		Name:  name,
		Value: *value,
	}, true, nil
}

func (s *Storage) ListCounters() (models.CountersList, error) {
	names, values, err := listMetrics[models.CounterValue](s, listCountersQuery)
	if err != nil {
		return nil, fmt.Errorf("list counters: %w", err)
	}
	counters, err := zipCounters(names, values)
	if err != nil {
		return nil, fmt.Errorf("zipping counters: %w", err)
	}
	return counters, nil
}

func (s *Storage) ReplaceCounters(val models.CountersList) error {
	names, values := unzipCounters(val)
	if err := replaceMetrics(s, clearCountersQuery, setCounterQuery, names, values); err != nil {
		return fmt.Errorf("replace gauges: %w", err)
	}
	return nil
}

func setMetric[Value any](s *Storage, query, name string, value Value) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("database not exists")
	}

	ctx, _ := context.WithTimeout(context.Background(), SqlOpTimeout)
	if _, err := s.db.ExecContext(ctx, query, name, value); err != nil {
		return fmt.Errorf("set metric: %w", err)
	}

	return nil
}

func findMetric[Value any](s *Storage, query, name string) (*Value, bool, error) {
	if s == nil || s.db == nil {
		return nil, false, fmt.Errorf("database not exists")
	}

	ctx, _ := context.WithTimeout(context.Background(), SqlOpTimeout)
	row, err := s.db.QueryRowContext(ctx, query, name)
	if err != nil {
		return nil, false, fmt.Errorf("find metric: %w", err)
	}

	var value Value
	if err := row.Scan(&value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		} else {
			return nil, false, fmt.Errorf("metric parsing: %w", err)
		}
	}

	return &value, true, nil
}

func listMetrics[Value any](s *Storage, query string) (names []string, values []Value, err error) {
	if s == nil || s.db == nil {
		return nil, nil, fmt.Errorf("database not exists")
	}

	// request rows
	ctx, _ := context.WithTimeout(context.Background(), SqlOpTimeout)
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, nil, fmt.Errorf("query rows: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			s.log.Error("metrics rows closing", zap.Error(err))
		}
	}()

	// read rows
	for rows.Next() {
		var (
			name  string
			value Value
		)
		if err := rows.Scan(&name, &value); err != nil {
			return nil, nil, err
		}
		names = append(names, name)
		values = append(values, value)
	}

	err = rows.Err()
	if err != nil {
		return nil, nil, fmt.Errorf("metrics rows reading: %w", err)
	}

	return names, values, nil
}

func replaceMetrics[Value any](s *Storage, clearQuery, addQuery string, names []string, values []Value) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("database not exists")
	}
	if len(names) != len(values) {
		return fmt.Errorf("insufficient names and values")
	}

	// start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("replace metrics: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			s.log.Error("replace metrics rollback", zap.Error(err))
		}
	}()
	ctx, _ := context.WithTimeout(context.Background(), SqlOpTimeout)

	// clear existed values
	if _, err = tx.ExecContext(ctx, clearQuery); err != nil {
		return fmt.Errorf("cleaning existed values: %w", err)
	}

	// add new values
	stmt, err := tx.PrepareContext(ctx, addQuery)
	if err != nil {
		return fmt.Errorf("prepare context: %w", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			s.log.Error("close add metric statement", zap.Error(err))
		}
	}()

	for i := range names {
		if _, err = stmt.ExecContext(ctx, names[i], values[i]); err != nil {
			return fmt.Errorf("add metric: %w", err)
		}
	}

	return nil
}

func zipGauges(names []string, values []models.GaugeValue) (models.GaugesList, error) {
	if len(names) != len(values) {
		return nil, fmt.Errorf("insufficient names and values")
	}
	gauges := make([]models.Gauge, len(names))
	for i := range names {
		gauges[i] = models.Gauge{
			Name:  names[i],
			Value: values[i],
		}
	}
	return gauges, nil
}

func zipCounters(names []string, values []models.CounterValue) (models.CountersList, error) {
	if len(names) != len(values) {
		return nil, fmt.Errorf("insufficient names and values")
	}
	counters := make([]models.Counter, len(names))
	for i := range names {
		counters[i] = models.Counter{
			Name:  names[i],
			Value: values[i],
		}
	}
	return counters, nil
}

func unzipGauges(gauges models.GaugesList) (names []string, values []models.GaugeValue) {
	names = make([]string, len(gauges))
	values = make([]models.GaugeValue, len(gauges))
	for i, gauge := range gauges {
		names[i] = gauge.Name
		values[i] = gauge.Value
	}
	return names, values
}

func unzipCounters(counters models.CountersList) (names []string, values []models.CounterValue) {
	names = make([]string, len(counters))
	values = make([]models.CounterValue, len(counters))
	for i, counter := range counters {
		names[i] = counter.Name
		values[i] = counter.Value
	}
	return names, values
}
