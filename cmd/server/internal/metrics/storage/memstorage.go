package storage

import (
	"sync"

	"github.com/stepkareserva/obsermon/internal/models"
)

var _ Storage = (*MemStorage)(nil)

type MemStorage struct {
	gauges   map[string]models.Gauge
	counters map[string]models.Counter
	lock     sync.Mutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]models.Gauge),
		counters: make(map[string]models.Counter),
	}
}

func (m *MemStorage) SetGauge(name string, val models.Gauge) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.gauges[name] = val
}

func (m *MemStorage) GetGauge(name string) (models.Gauge, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()

	val, exists := m.gauges[name]
	return val, exists
}

func (m *MemStorage) ListGauges() models.Names {
	m.lock.Lock()
	defer m.lock.Unlock()

	return listKeys(m.gauges)
}

func (m *MemStorage) SetCounter(name string, val models.Counter) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.counters[name] = val
}

func (m *MemStorage) GetCounter(name string) (models.Counter, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()

	val, exists := m.counters[name]
	return val, exists
}

func (m *MemStorage) ListCounters() models.Names {
	m.lock.Lock()
	defer m.lock.Unlock()

	return listKeys(m.counters)
}

func listKeys[T any](m map[string]T) models.Names {
	keys := make(models.Names)
	for key := range m {
		keys[key] = struct{}{}
	}
	return keys
}
