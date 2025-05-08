package handlers

import (
	"context"

	"github.com/stepkareserva/obsermon/internal/models"
)

type GaugesService interface {
	UpdateGauge(ctx context.Context, val models.Gauge) (*models.Gauge, error)
	FindGauge(ctx context.Context, name string) (*models.Gauge, bool, error)
	ListGauges(ctx context.Context) (models.GaugesList, error)
}

type CountersService interface {
	UpdateCounter(ctx context.Context, val models.Counter) (*models.Counter, error)
	FindCounter(ctx context.Context, name string) (*models.Counter, bool, error)
	ListCounters(ctx context.Context) (models.CountersList, error)
}

type MetricsService interface {
	UpdateMetric(ctx context.Context, val models.Metric) (*models.Metric, error)
	FindMetric(ctx context.Context, t models.MetricType, name string) (*models.Metric, bool, error)
	UpdateMetrics(ctx context.Context, vals models.Metrics) (models.Metrics, error)
}

//go:generate mockgen -source=$GOFILE -destination=../../mocks/mock_service.go -package=mocks

type Service interface {
	GaugesService
	CountersService
	MetricsService
}
