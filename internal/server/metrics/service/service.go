package service

import (
	"fmt"
	"sort"

	"github.com/stepkareserva/obsermon/internal/models"
	"github.com/stepkareserva/obsermon/internal/server/handlers"
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

func (s *Service) Close() error {
	return nil
}

func (s *Service) UpdateGauge(val models.Gauge) (*models.Gauge, error) {
	if err := s.checkValidity(); err != nil {
		return nil, err
	}

	if err := s.storage.SetGauge(val); err != nil {
		return nil, err
	}

	return &val, nil
}

func (s *Service) FindGauge(name string) (*models.Gauge, bool, error) {
	if err := s.checkValidity(); err != nil {
		return nil, false, err
	}

	return s.storage.FindGauge(name)
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

func (s *Service) UpdateCounter(val models.Counter) (*models.Counter, error) {
	if err := s.checkValidity(); err != nil {
		return nil, err
	}

	current, exists, err := s.storage.FindCounter(val.Name)
	if err != nil {
		return nil, err
	}

	if !exists {
		if err := s.storage.SetCounter(val); err != nil {
			return nil, err
		}
		return &val, nil
	}

	if err = current.Value.Update(val.Value); err != nil {
		return nil, err
	}
	if err = s.storage.SetCounter(*current); err != nil {
		return nil, err
	}
	return &val, nil
}

func (s *Service) FindCounter(name string) (*models.Counter, bool, error) {
	if err := s.checkValidity(); err != nil {
		return nil, false, err
	}

	return s.storage.FindCounter(name)
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

func (s *Service) UpdateMetric(val models.Metrics) (*models.Metrics, error) {
	if err := s.checkValidity(); err != nil {
		return nil, err
	}

	switch val.MType {
	case models.MetricTypeCounter:
		counter, err := val.Counter()
		if err != nil {
			return nil, err
		}
		updated, err := s.UpdateCounter(*counter)
		if err != nil {
			return nil, err
		}
		updatedMetric := models.CounterMetric(*updated)
		return &updatedMetric, nil
	case models.MetricTypeGauge:
		gauge, err := val.Gauge()
		if err != nil {
			return nil, err
		}
		updated, err := s.UpdateGauge(*gauge)
		if err != nil {
			return nil, err
		}
		updatedMetric := models.GaugeMetric(*updated)
		return &updatedMetric, nil
	default:
		return nil, fmt.Errorf("unknown metric type")
	}
}

func (s *Service) FindMetric(t models.MetricType, name string) (*models.Metrics, bool, error) {
	if err := s.checkValidity(); err != nil {
		return nil, false, err
	}

	switch t {
	case models.MetricTypeCounter:
		c, exists, err := s.FindCounter(name)
		if err != nil || !exists {
			return nil, exists, err
		}
		m := models.CounterMetric(*c)
		return &m, true, nil
	case models.MetricTypeGauge:
		g, exists, err := s.FindGauge(name)
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
