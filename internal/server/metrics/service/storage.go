package service

import (
	"context"

	"github.com/stepkareserva/obsermon/internal/models"
)

type GaugeStorage interface {
	SetGauge(ctx context.Context, val models.Gauge) error
	SetGauges(ctx context.Context, val models.GaugesList) error
	FindGauge(ctx context.Context, name string) (*models.Gauge, bool, error)
	ListGauges(ctx context.Context) (models.GaugesList, error)
	ReplaceGauges(ctx context.Context, val models.GaugesList) error
}

type CounterOp = func(val *models.CounterValue) error

type CounterStorage interface {
	UpdateCounter(ctx context.Context, val models.Counter) (*models.Counter, error)
	UpdateCounters(ctx context.Context, vals models.CountersList) (models.CountersList, error)
	FindCounter(ctx context.Context, name string) (*models.Counter, bool, error)
	ListCounters(ctx context.Context) (models.CountersList, error)
	ReplaceCounters(ctx context.Context, val models.CountersList) error
}

type Pingable interface {
	Ping(ctx context.Context) error
}

//go:generate mockgen -source=$GOFILE -destination=../../mocks/mock_storage.go -package=mocks

type Storage interface {
	GaugeStorage
	CounterStorage
}
