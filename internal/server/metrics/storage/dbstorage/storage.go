package dbstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/stepkareserva/obsermon/internal/models"
	"github.com/stepkareserva/obsermon/internal/server/interfaces/database"
	"github.com/stepkareserva/obsermon/internal/server/metrics/service"
	"go.uber.org/zap"
)

type Storage struct {
	db  database.Database
	log *zap.Logger
}

var _ service.Storage = (*Storage)(nil)

const (
	CountersTable = "counters"
	GaugesTable   = "gauges"

	NameColumn  = "name"
	ValueColumn = "value"

	SQLOpTimeout = 15 * time.Second
)

func New(db database.Database, log *zap.Logger) (*Storage, error) {
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

	if _, err := s.db.Exec(ctx, createCountersQuery); err != nil {
		return fmt.Errorf("counters table creation: %w", err)
	}

	// create gauges table
	if _, err := s.db.Exec(ctx, createGaugesQuery); err != nil {
		return fmt.Errorf("gauges table creation: %w", err)
	}

	return nil
}

func (s *Storage) SetGauge(ctx context.Context, val models.Gauge) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("database not exists")
	}

	if _, err := s.db.Exec(ctx, setGaugeQuery, val.Name, val.Value); err != nil {
		return fmt.Errorf("set gauge: %w", err)
	}

	return nil
}

func (s *Storage) SetGauges(ctx context.Context, vals models.GaugesList) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("database not exists")
	}

	txfn := func(ctx context.Context, tx *sql.Tx) (err error) {
		// insert or update gauges
		setGaugeStmt, err := tx.PrepareContext(ctx, setGaugeQuery)
		if err != nil {
			return fmt.Errorf("prepare set gauges context: %w", err)
		}
		defer func() {
			if err = setGaugeStmt.Close(); err != nil {
				err = fmt.Errorf("close set gauges statement: %w", err)
			}
		}()

		for _, val := range vals {
			if _, err = setGaugeStmt.ExecContext(ctx, val.Name, val.Value); err != nil {
				return fmt.Errorf("set gauge: %w", err)
			}
		}
		return nil
	}

	if err := s.db.ExecTxFn(ctx, txfn); err != nil {
		return fmt.Errorf("set gauges: %w", err)
	}

	return nil
}
func (s *Storage) FindGauge(ctx context.Context, name string) (*models.Gauge, bool, error) {
	value, exists, err := findMetric[models.GaugeValue](ctx, s, findGaugeQuery, name)
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

func (s *Storage) ListGauges(ctx context.Context) (models.GaugesList, error) {
	names, values, err := listMetrics[models.GaugeValue](ctx, s, listGaugesQuery)
	if err != nil {
		return nil, fmt.Errorf("list gauges: %w", err)
	}
	gauges, err := zipGauges(names, values)
	if err != nil {
		return nil, fmt.Errorf("zipping gauges: %w", err)
	}
	return gauges, nil
}

func (s *Storage) ReplaceGauges(ctx context.Context, val models.GaugesList) error {
	names, values := unzipGauges(val)
	if err := replaceMetrics(ctx, s, clearGaugeQuery, insertGaugeQuery, names, values); err != nil {
		return fmt.Errorf("replace gauges: %w", err)
	}
	return nil
}

func (s *Storage) UpdateCounter(ctx context.Context, val models.Counter) (*models.Counter, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("database not exists")
	}

	var updatedVal models.Counter
	txfn := func(ctx context.Context, tx *sql.Tx) error {
		// get current counter value
		var currentVal models.CounterValue
		err := tx.QueryRowContext(ctx, selectCounterForUpdateQuery, val.Name).Scan(&currentVal)

		// update or insert updated counter value
		switch {
		case err == sql.ErrNoRows:
			if _, err = tx.ExecContext(ctx, insertCounterQuery, val.Name, val.Value); err != nil {
				return fmt.Errorf("insert counter: %w", err)
			}
			updatedVal = val
			return nil
		case err != nil:
			return fmt.Errorf("query counter: %w", err)
		default:
			if err = val.Value.Update(currentVal); err != nil {
				return fmt.Errorf("update counter value: %w", err)
			}
			if _, err := tx.ExecContext(ctx, updateCounterQuery, val.Name, val.Value); err != nil {
				return fmt.Errorf("update counter: %w", err)
			}
			updatedVal = val
			return nil
		}
	}

	if err := s.db.ExecTxFn(ctx, txfn); err != nil {
		return nil, fmt.Errorf("update counter: %w", err)
	}

	return &updatedVal, nil
}

func (s *Storage) UpdateCounters(ctx context.Context, vals models.CountersList) (models.CountersList, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("database not exists")
	}

	var updatedVals models.CountersList
	txfn := func(ctx context.Context, tx *sql.Tx) (err error) {
		// lock counter statement
		selectCounterStmt, err := tx.PrepareContext(ctx, selectCounterForUpdateQuery)
		if err != nil {
			return fmt.Errorf("prepare select counter statement: %w", err)
		}
		defer func() {
			if closeErr := selectCounterStmt.Close(); closeErr != nil {
				err = fmt.Errorf("close select counter statement: %w", closeErr)
			}
		}()

		// insert counter statement
		insertCounterStmt, err := tx.PrepareContext(ctx, insertCounterQuery)
		if err != nil {
			return fmt.Errorf("prepare insert counter statement: %w", err)
		}
		defer func() {
			if closeErr := insertCounterStmt.Close(); closeErr != nil {
				err = fmt.Errorf("close insert counter statement: %w", closeErr)
			}
		}()

		// update counter statemnt
		updateCounterStmt, err := tx.PrepareContext(ctx, updateCounterQuery)
		if err != nil {
			return fmt.Errorf("prepare update counter statement: %w", err)
		}
		defer func() {
			if closeErr := updateCounterStmt.Close(); closeErr != nil {
				err = fmt.Errorf("close update counter statement: %w", closeErr)
			}
		}()

		for _, val := range vals {
			// get counter value
			var currentVal models.CounterValue
			err = selectCounterStmt.QueryRowContext(ctx, val.Name).Scan(&currentVal)

			// update or insert updated counter value
			switch {
			case err == sql.ErrNoRows:
				if _, err = insertCounterStmt.ExecContext(ctx, val.Name, val.Value); err != nil {
					return fmt.Errorf("insert counter: %w", err)
				}
				updatedVals = append(updatedVals, val)
			case err != nil:
				return fmt.Errorf("query counter: %w", err)
			default:
				if err = val.Value.Update(currentVal); err != nil {
					return fmt.Errorf("update counter value: %w", err)
				}
				if _, err := updateCounterStmt.ExecContext(ctx, val.Name, val.Value); err != nil {
					return fmt.Errorf("update counter: %w", err)
				}
				updatedVals = append(updatedVals, val)
			}
		}

		return nil
	}

	if err := s.db.ExecTxFn(ctx, txfn); err != nil {
		return nil, fmt.Errorf("update counters: %w", err)
	}

	return updatedVals, nil
}

func (s *Storage) FindCounter(ctx context.Context, name string) (*models.Counter, bool, error) {
	value, exists, err := findMetric[models.CounterValue](ctx, s, findCounterQuery, name)
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

func (s *Storage) ListCounters(ctx context.Context) (models.CountersList, error) {
	names, values, err := listMetrics[models.CounterValue](ctx, s, listCountersQuery)
	if err != nil {
		return nil, fmt.Errorf("list counters: %w", err)
	}
	counters, err := zipCounters(names, values)
	if err != nil {
		return nil, fmt.Errorf("zipping counters: %w", err)
	}
	return counters, nil
}

func (s *Storage) ReplaceCounters(ctx context.Context, val models.CountersList) error {
	names, values := unzipCounters(val)
	if err := replaceMetrics(ctx, s, clearCountersQuery, insertCounterQuery, names, values); err != nil {
		return fmt.Errorf("replace counters: %w", err)
	}
	return nil
}

func findMetric[Value any](ctx context.Context, s *Storage, query, name string) (*Value, bool, error) {
	if s == nil || s.db == nil {
		return nil, false, fmt.Errorf("database not exists")
	}

	row, err := s.db.QueryRow(ctx, query, name)
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

func listMetrics[Value any](ctx context.Context, s *Storage, query string) (names []string, values []Value, err error) {
	if s == nil || s.db == nil {
		return nil, nil, fmt.Errorf("database not exists")
	}

	rows, err := s.db.Query(ctx, query)
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

func replaceMetrics[Value any](ctx context.Context, s *Storage, clearQuery, insertQuery string, names []string, values []Value) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("database not exists")
	}
	if len(names) != len(values) {
		return fmt.Errorf("insufficient names and values")
	}

	txfn := func(ctx context.Context, tx *sql.Tx) error {
		// clear existed values
		if _, err := tx.ExecContext(ctx, clearQuery); err != nil {
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

	if err := s.db.ExecTxFn(ctx, txfn); err != nil {
		return fmt.Errorf("replace metrics: %w", err)
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
