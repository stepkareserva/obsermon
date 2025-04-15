package handlers

import "github.com/stepkareserva/obsermon/internal/models"

type GaugesService interface {
	UpdateGauge(val models.Gauge) error
	GetGauge(name string) (*models.Gauge, bool, error)
	ListGauges() (models.GaugesList, error)
}

type CountersService interface {
	UpdateCounter(val models.Counter) error
	GetCounter(name string) (*models.Counter, bool, error)
	ListCounters() (models.CountersList, error)
}

type MetricsService interface {
	UpdateMetric(val models.Metrics) error
	GetMetric(t models.MetricType, name string) (*models.Metrics, bool, error)
}

//go:generate mockgen -source=$GOFILE -destination=../../mocks/mock_service.go -package=mocks

type Service interface {
	GaugesService
	CountersService
	MetricsService
}
