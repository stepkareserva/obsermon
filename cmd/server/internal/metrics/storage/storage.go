package storage

import "github.com/stepkareserva/obsermon/internal/models"

type Storage interface {
	SetGauge(name string, val models.Gauge)
	GetGauge(name string) (models.Gauge, bool)
	ListGauges() models.Names

	SetCounter(name string, val models.Counter)
	GetCounter(name string) (models.Counter, bool)
	ListCounters() models.Names
}
