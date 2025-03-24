package storage

import "github.com/stepkareserva/obsermon/internal/models"

type Storage interface {
	SetGauge(val models.Gauge) error
	GetGauge(name string) (*models.Gauge, bool, error)
	ListGauges() (models.GaugesList, error)

	SetCounter(val models.Counter) error
	GetCounter(name string) (*models.Counter, bool, error)
	ListCounters() (models.CountersList, error)
}
