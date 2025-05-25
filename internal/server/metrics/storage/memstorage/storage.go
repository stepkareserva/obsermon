package memstorage

import (
	"context"
	"fmt"
	"sync"

	"github.com/stepkareserva/obsermon/internal/models"
	"github.com/stepkareserva/obsermon/internal/server/metrics/service"
)

var _ service.Storage = (*Storage)(nil)

type Storage struct {
	gauges   models.GaugesMap
	counters models.CountersMap
	lock     sync.RWMutex
}

func New() *Storage {
	return &Storage{
		gauges:   make(models.GaugesMap),
		counters: make(models.CountersMap),
	}
}

func (m *Storage) SetGauge(ctx context.Context, val models.Gauge) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.gauges[val.Name] = val.Value
	return nil
}

func (m *Storage) SetGauges(ctx context.Context, vals models.GaugesList) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, val := range vals {
		m.gauges[val.Name] = val.Value
	}
	return nil
}
func (m *Storage) FindGauge(ctx context.Context, name string) (*models.Gauge, bool, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	val, exists := m.gauges[name]
	return &models.Gauge{Name: name, Value: val}, exists, nil
}

func (m *Storage) ListGauges(ctx context.Context) (models.GaugesList, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.gauges.List(), nil
}

func (m *Storage) ReplaceGauges(ctx context.Context, val models.GaugesList) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.gauges = val.Map()
	return nil
}

func (m *Storage) UpdateCounter(ctx context.Context, val models.Counter) (*models.Counter, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	counter, exists := m.counters[val.Name]
	if exists {
		if err := val.Value.Update(counter); err != nil {
			return nil, fmt.Errorf("update counter: %v", err)
		}
	}
	m.counters[val.Name] = val.Value

	return &val, nil
}

func (m *Storage) UpdateCounters(ctx context.Context, vals models.CountersList) (models.CountersList, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, val := range vals {
		counter, exists := m.counters[val.Name]
		if exists {
			if err := val.Value.Update(counter); err != nil {
				return nil, fmt.Errorf("update counters: %v", err)
			}
		}
		m.counters[val.Name] = val.Value
	}

	return vals, nil
}

func (m *Storage) FindCounter(ctx context.Context, name string) (*models.Counter, bool, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	val, exists := m.counters[name]
	return &models.Counter{Name: name, Value: val}, exists, nil
}

func (m *Storage) ListCounters(ctx context.Context) (models.CountersList, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.counters.List(), nil
}

func (m *Storage) ReplaceCounters(ctx context.Context, val models.CountersList) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.counters = val.Map()
	return nil
}
