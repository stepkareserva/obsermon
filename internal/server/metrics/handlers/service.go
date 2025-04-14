package handlers

import "github.com/stepkareserva/obsermon/internal/models"

type GaugesService interface {
	UpdateGauge(val models.Gauge) error
	GetGauge(name string) (*models.Gauge, bool, error)
	ListGauges() ([]models.Gauge, error)
}

type CountersService interface {
	UpdateCounter(val models.Counter) error
	GetCounter(name string) (*models.Counter, bool, error)
	ListCounters() ([]models.Counter, error)
}

//go:generate mockgen -source=$GOFILE -destination=../../mocks/mock_service.go -package=mocks

type Service interface {
	GaugesService
	CountersService
}
