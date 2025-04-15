package service

import (
	"fmt"
	"sort"

	"github.com/stepkareserva/obsermon/internal/models"
	"github.com/stepkareserva/obsermon/internal/server/metrics/handlers"
)

type Service struct {
	storage Storage
}

var _ handlers.Service = (*Service)(nil)

func New(storage Storage) (*Service, error) {
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

func (s *Service) ListGauges() (models.GaugesList, error) {
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

func (s *Service) ReplaceGauges(val models.GaugesList) error {
	if err := s.checkValidity(); err != nil {
		return err
	}
	return s.storage.ReplaceGauges(val)
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

func (s *Service) ListCounters() (models.CountersList, error) {
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

func (s *Service) ReplaceCounters(val models.CountersList) error {
	if err := s.checkValidity(); err != nil {
		return err
	}
	return s.storage.ReplaceCounters(val)
}

func (s *Service) UpdateMetric(val models.Metrics) error {
	if err := s.checkValidity(); err != nil {
		return err
	}

	switch val.MType {
	case models.MetricTypeCounter:
		counter, err := val.Counter()
		if err != nil {
			return err
		}
		return s.UpdateCounter(*counter)
	case models.MetricTypeGauge:
		gauge, err := val.Gauge()
		if err != nil {
			return err
		}
		return s.UpdateGauge(*gauge)
	default:
		return fmt.Errorf("unknown metric type")
	}
}

func (s *Service) GetMetric(t models.MetricType, name string) (*models.Metrics, bool, error) {
	if err := s.checkValidity(); err != nil {
		return nil, false, err
	}

	switch t {
	case models.MetricTypeCounter:
		c, exists, err := s.GetCounter(name)
		if err != nil || !exists {
			return nil, exists, err
		}
		m := models.CounterMetric(*c)
		return &m, true, nil
	case models.MetricTypeGauge:
		g, exists, err := s.GetGauge(name)
		if err != nil || !exists {
			return nil, exists, err
		}
		m := models.GaugeMetric(*g)
		return &m, true, nil
	default:
		return nil, false, fmt.Errorf("unknown metric type")
	}
}

func (s *Service) checkValidity() error {
	if s == nil || s.storage == nil {
		return fmt.Errorf("Service not exists")
	}
	return nil
}
