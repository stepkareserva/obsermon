package handlers

import "github.com/stepkareserva/obsermon/internal/models"

type GaugesService interface {
	UpdateGauge(val models.Gauge) (*models.Gauge, error)
	FindGauge(name string) (*models.Gauge, bool, error)
	ListGauges() (models.GaugesList, error)
}

type CountersService interface {
	UpdateCounter(val models.Counter) (*models.Counter, error)
	FindCounter(name string) (*models.Counter, bool, error)
	ListCounters() (models.CountersList, error)
}

type MetricsService interface {
	UpdateMetric(val models.Metrics) (*models.Metrics, error)
	FindMetric(t models.MetricType, name string) (*models.Metrics, bool, error)
	UpdateMetrics(vals []models.Metrics) ([]models.Metrics, error)
}

//go:generate mockgen -source=$GOFILE -destination=../../mocks/mock_service.go -package=mocks

type Service interface {
	GaugesService
	CountersService
	MetricsService
}
