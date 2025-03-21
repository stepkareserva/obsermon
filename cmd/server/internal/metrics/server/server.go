package server

import (
	"fmt"

	"github.com/stepkareserva/obsermon/cmd/server/internal/metrics/storage"
	"github.com/stepkareserva/obsermon/internal/models"
)

type Server struct {
	storage storage.Storage
}

func NewServer(storage storage.Storage) (*Server, error) {
	if storage == nil {
		return nil, fmt.Errorf("metrics storage is nil")
	}
	return &Server{storage: storage}, nil
}

func (s *Server) UpdateGauge(name string, val models.Gauge) error {
	if err := s.checkValidity(); err != nil {
		return err
	}

	s.storage.SetGauge(name, val)
	return nil
}

func (s *Server) GetGauge(name string) (models.Gauge, error) {
	if err := s.checkValidity(); err != nil {
		return 0, err
	}

	val, exists := s.storage.GetGauge(name)
	if !exists {
		return 0, fmt.Errorf("gauge does not exist")
	}

	return val, nil
}

func (s *Server) ListGauges() (models.Names, error) {
	if err := s.checkValidity(); err != nil {
		return nil, err
	}
	return s.storage.ListGauges(), nil
}

func (s *Server) UpdateCounter(name string, val models.Counter) error {
	if err := s.checkValidity(); err != nil {
		return err
	}

	current, exists := s.storage.GetCounter(name)
	if !exists {
		s.storage.SetCounter(name, val)
		return nil
	}

	err := current.Update(val)
	if err != nil {
		return err
	}

	s.storage.SetCounter(name, current)
	return nil
}

func (s *Server) GetCounter(name string) (models.Counter, error) {
	if err := s.checkValidity(); err != nil {
		return 0, err
	}

	value, exists := s.storage.GetCounter(name)
	if !exists {
		return value, fmt.Errorf("counter does not exist")
	}

	return value, nil
}

func (s *Server) ListCounters() (models.Names, error) {
	if err := s.checkValidity(); err != nil {
		return nil, err
	}

	return s.storage.ListCounters(), nil
}

func (s *Server) checkValidity() error {
	if s == nil || s.storage == nil {
		return fmt.Errorf("Server not exists")
	}
	return nil
}
