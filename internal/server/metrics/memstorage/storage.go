package memstorage

import (
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

func (m *Storage) SetGauge(val models.Gauge) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.gauges[val.Name] = val.Value
	return nil
}

func (m *Storage) FindGauge(name string) (*models.Gauge, bool, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	val, exists := m.gauges[name]
	return &models.Gauge{Name: name, Value: val}, exists, nil
}

func (m *Storage) ListGauges() (models.GaugesList, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.gauges.List(), nil
}

func (m *Storage) ReplaceGauges(val models.GaugesList) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.gauges = val.Map()
	return nil
}

func (m *Storage) SetCounter(val models.Counter) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.counters[val.Name] = val.Value
	return nil
}

func (m *Storage) FindCounter(name string) (*models.Counter, bool, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	val, exists := m.counters[name]
	return &models.Counter{Name: name, Value: val}, exists, nil
}

func (m *Storage) ListCounters() (models.CountersList, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.counters.List(), nil
}

func (m *Storage) ReplaceCounters(val models.CountersList) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.counters = val.Map()
	return nil
}
