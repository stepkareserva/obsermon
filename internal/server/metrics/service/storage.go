package service

import "github.com/stepkareserva/obsermon/internal/models"

type GaugeStorage interface {
	SetGauge(val models.Gauge) error
	GetGauge(name string) (*models.Gauge, bool, error)
	ListGauges() (models.GaugesList, error)
	ReplaceGauges(val models.GaugesList) error
}

type CounterStorage interface {
	SetCounter(val models.Counter) error
	GetCounter(name string) (*models.Counter, bool, error)
	ListCounters() (models.CountersList, error)
	ReplaceCounters(val models.CountersList) error
}

//go:generate mockgen -source=$GOFILE -destination=../../mocks/mock_storage.go -package=mocks

type Storage interface {
	GaugeStorage
	CounterStorage
}
