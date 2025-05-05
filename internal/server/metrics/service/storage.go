package service

import "github.com/stepkareserva/obsermon/internal/models"

type GaugeStorage interface {
	SetGauge(val models.Gauge) error
	SetGauges(val models.GaugesList) error
	FindGauge(name string) (*models.Gauge, bool, error)
	ListGauges() (models.GaugesList, error)
	ReplaceGauges(val models.GaugesList) error
}

type CounterOp = func(val *models.CounterValue) error

type CounterStorage interface {
	UpdateCounter(val models.Counter) (*models.Counter, error)
	FindCounter(name string) (*models.Counter, bool, error)
	ListCounters() (models.CountersList, error)
	ReplaceCounters(val models.CountersList) error
}

//go:generate mockgen -source=$GOFILE -destination=../../mocks/mock_storage.go -package=mocks

type Storage interface {
	GaugeStorage
	CounterStorage
}
