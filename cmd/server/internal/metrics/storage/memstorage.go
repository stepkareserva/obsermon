package storage

import (
	"sync"

	"github.com/stepkareserva/obsermon/internal/models"
)

var _ Storage = (*MemStorage)(nil)

type MemStorage struct {
	gauges   models.GaugesMap
	counters models.CountersMap
	lock     sync.RWMutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(models.GaugesMap),
		counters: make(models.CountersMap),
	}
}

func (m *MemStorage) SetGauge(val models.Gauge) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.gauges[val.Name] = val.Value
	return nil
}

func (m *MemStorage) GetGauge(name string) (*models.Gauge, bool, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	val, exists := m.gauges[name]
	return &models.Gauge{Name: name, Value: val}, exists, nil
}

func (m *MemStorage) ListGauges() (models.GaugesList, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.gauges.List(), nil
}

func (m *MemStorage) SetCounter(val models.Counter) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.counters[val.Name] = val.Value
	return nil
}

func (m *MemStorage) GetCounter(name string) (*models.Counter, bool, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	val, exists := m.counters[name]
	return &models.Counter{Name: name, Value: val}, exists, nil
}

func (m *MemStorage) ListCounters() (models.CountersList, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.counters.List(), nil
}
