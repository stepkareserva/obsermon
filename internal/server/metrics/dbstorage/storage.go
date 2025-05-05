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

	SQLOpTimeout = 5 * time.Second
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

	insertCounterQuery = queryReplacer.Replace(`
		INSERT
			INTO {counters} ({name}, {value})
			VALUES ($1, $2)
		`)

	updateCounterQuery = queryReplacer.Replace(`
		UPDATE {counters}
			SET {value} = $2
			WHERE {name} = $1
		`)

	findCounterQuery = queryReplacer.Replace(`
		SELECT {value}
			FROM {counters}
			WHERE {name} = $1
		`)

	listCountersQuery = queryReplacer.Replace(`
		SELECT {name}, {value}
			FROM {counters}
		`)

	selectCounterForUpdateQuery = queryReplacer.Replace(`
		SELECT {value}
			FROM {counters}
			WHERE {name} = $1
		FOR UPDATE
		`)

	clearCountersQuery = queryReplacer.Replace(`
		DELETE FROM {counters}
	`)

	setGaugeQuery = queryReplacer.Replace(`
		INSERT
			INTO {gauges} ({name}, {value})
			VALUES ($1, $2)
		ON CONFLICT ({name})
			DO UPDATE SET {value} = EXCLUDED.{value}
		`)

	insertGaugeQuery = queryReplacer.Replace(`
		INSERT
			INTO {gauges} ({name}, {value})
			VALUES ($1, $2)
		`)

	findGaugeQuery = queryReplacer.Replace(`
		SELECT {value}
			FROM {gauges}
			WHERE {name} = $1
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
	ctx, cancel := context.WithTimeout(context.Background(), SQLOpTimeout)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, createCountersQuery); err != nil {
		return fmt.Errorf("counters table creation: %w", err)
	}

	// create gauges table
	if _, err := s.db.ExecContext(ctx, createGaugesQuery); err != nil {
		return fmt.Errorf("gauges table creation: %w", err)
	}

	return nil
}

func (s *Storage) SetGauge(val models.Gauge) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("database not exists")
	}

	ctx, cancel := context.WithTimeout(context.Background(), SQLOpTimeout)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, setGaugeQuery, val.Name, val.Value); err != nil {
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
	if err := replaceMetrics(s, clearGaugeQuery, insertGaugeQuery, names, values); err != nil {
		return fmt.Errorf("replace gauges: %w", err)
	}
	return nil
}

func (s *Storage) UpdateCounter(val models.Counter) (*models.Counter, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("database not exists")
	}

	ctx, cancel := context.WithTimeout(context.Background(), SQLOpTimeout)
	defer cancel()

	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer s.txClosingFn(tx, err)

	// get current counter value
	var currentVal models.CounterValue
	err = tx.QueryRowContext(ctx, selectCounterForUpdateQuery, val.Name).Scan(&currentVal)

	// update or insert updated counter value
	switch {
	case err == sql.ErrNoRows:
		if _, err = tx.ExecContext(ctx, insertCounterQuery, val.Name, val.Value); err != nil {
			return nil, fmt.Errorf("insert counter: %w", err)
		}
		return &val, nil
	case err != nil:
		return nil, fmt.Errorf("query counter: %w", err)
	default:
		if err = val.Value.Update(currentVal); err != nil {
			return nil, fmt.Errorf("update counter value: %w", err)
		}
		if _, err := tx.ExecContext(ctx, updateCounterQuery, val.Name, val.Value); err != nil {
			return nil, fmt.Errorf("update counter: %w", err)
		}
		return &val, nil
	}
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
	if err := replaceMetrics(s, clearCountersQuery, insertCounterQuery, names, values); err != nil {
		return fmt.Errorf("replace counters: %w", err)
	}
	return nil
}

func (s *Storage) txClosingFn(tx *sql.Tx, err error) {
	if err != nil {
		err = tx.Rollback()
	} else {
		err = tx.Commit()
	}
	if err != nil {
		s.log.Error("update counter transaction", zap.Error(err))
	}
}

func findMetric[Value any](s *Storage, query, name string) (*Value, bool, error) {
	if s == nil || s.db == nil {
		return nil, false, fmt.Errorf("database not exists")
	}

	ctx, cancel := context.WithTimeout(context.Background(), SQLOpTimeout)
	defer cancel()

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

	ctx, cancel := context.WithTimeout(context.Background(), SQLOpTimeout)
	defer cancel()

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

func replaceMetrics[Value any](s *Storage, clearQuery, insertQuery string, names []string, values []Value) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("database not exists")
	}
	if len(names) != len(values) {
		return fmt.Errorf("insufficient names and values")
	}

	ctx, cancel := context.WithTimeout(context.Background(), SQLOpTimeout)
	defer cancel()

	// start transaction
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("replace metrics: %w", err)
	}
	defer s.txClosingFn(tx, err)

	// clear existed values
	if _, err = tx.ExecContext(ctx, clearQuery); err != nil {
		return fmt.Errorf("cleaning existed values: %w", err)
	}

	// add new values
	stmt, err := tx.PrepareContext(ctx, insertQuery)
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
