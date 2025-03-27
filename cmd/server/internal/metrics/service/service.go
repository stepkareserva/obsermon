package service

import (
	"fmt"
	"sort"

	"github.com/stepkareserva/obsermon/cmd/server/internal/metrics/storage"
	"github.com/stepkareserva/obsermon/internal/models"
)

type Service struct {
	storage storage.Storage
}

func NewService(storage storage.Storage) (*Service, error) {
	if storage == nil {
		return nil, fmt.Errorf("metrics storage is nil")
	}
	return &Service{storage: storage}, nil
}

func (s *Service) UpdateGauge(val models.Gauge) error {
	if err := s.checkValidity(); err != nil {
		return err
	}

	return s.storage.SetGauge(val)
}

func (s *Service) GetGauge(name string) (*models.Gauge, bool, error) {
	if err := s.checkValidity(); err != nil {
		return nil, false, err
	}

	return s.storage.GetGauge(name)
}

func (s *Service) ListGauges() ([]models.Gauge, error) {
	if err := s.checkValidity(); err != nil {
		return nil, err
	}
	gauges, err := s.storage.ListGauges()
	if err != nil {
		return nil, err
	}

	sort.SliceStable(gauges, func(i, j int) bool {
		return gauges[i].Name < gauges[j].Name
	})

	return gauges, nil
}

func (s *Service) UpdateCounter(val models.Counter) error {
	if err := s.checkValidity(); err != nil {
		return err
	}

	current, exists, err := s.storage.GetCounter(val.Name)
	if err != nil {
		return err
	}
	if !exists {
		return s.storage.SetCounter(val)
	}

	if err = current.Value.Update(val.Value); err != nil {
		return err
	}

	return s.storage.SetCounter(*current)
}

func (s *Service) GetCounter(name string) (*models.Counter, bool, error) {
	if err := s.checkValidity(); err != nil {
		return nil, false, err
	}

	return s.storage.GetCounter(name)
}

func (s *Service) ListCounters() ([]models.Counter, error) {
	if err := s.checkValidity(); err != nil {
		return nil, err
	}

	counters, err := s.storage.ListCounters()
	if err != nil {
		return nil, err
	}

	sort.SliceStable(counters, func(i, j int) bool {
		return counters[i].Name < counters[j].Name
	})

	return counters, nil
}

func (s *Service) checkValidity() error {
	if s == nil || s.storage == nil {
		return fmt.Errorf("Service not exists")
	}
	return nil
}
