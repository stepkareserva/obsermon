package persistence

import (
	"fmt"

	"github.com/stepkareserva/obsermon/internal/models"
	"github.com/stepkareserva/obsermon/internal/server/metrics/service"
)

type State struct {
	Counters []models.Counter `json:"counters"`
	Gauges   []models.Gauge   `json:"gauge"`
}

func (s *State) Export(storage service.Storage) error {
	if s == nil {
		return fmt.Errorf("state not exists")
	}
	if storage == nil {
		return fmt.Errorf("storage not exists")
	}
	if err := storage.ReplaceCounters(s.Counters); err != nil {
		return fmt.Errorf("replacing storage counters: %w", err)
	}
	if err := storage.ReplaceGauges(s.Gauges); err != nil {
		return fmt.Errorf("replacing storage gauges: %w", err)
	}
	return nil
}

func (s *State) Import(storage service.Storage) error {
	if s == nil {
		return fmt.Errorf("state not exists")
	}
	if storage == nil {
		return fmt.Errorf("storage not exists")
	}
	if err := storage.ReplaceCounters(s.Counters); err != nil {
		return fmt.Errorf("replacing storage counters: %w", err)
	}
	if err := storage.ReplaceGauges(s.Gauges); err != nil {
		return fmt.Errorf("replacing storage gauges: %w", err)
	}
	return nil
}
