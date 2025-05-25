package persistence

import (
	"context"
	"fmt"

	"github.com/stepkareserva/obsermon/internal/models"
	"github.com/stepkareserva/obsermon/internal/server/metrics/service"
)

type State struct {
	Counters []models.Counter `json:"counters"`
	Gauges   []models.Gauge   `json:"gauge"`
}

func (s *State) Export(ctx context.Context, storage service.Storage) error {
	if s == nil {
		return fmt.Errorf("state not exists")
	}
	if storage == nil {
		return fmt.Errorf("storage not exists")
	}
	if err := storage.ReplaceCounters(ctx, s.Counters); err != nil {
		return fmt.Errorf("replacing storage counters: %v", err)
	}
	if err := storage.ReplaceGauges(ctx, s.Gauges); err != nil {
		return fmt.Errorf("replacing storage gauges: %v", err)
	}
	return nil
}

func (s *State) Import(ctx context.Context, storage service.Storage) error {
	if s == nil {
		return fmt.Errorf("state not exists")
	}
	if storage == nil {
		return fmt.Errorf("storage not exists")
	}
	var err error
	s.Counters, err = storage.ListCounters(ctx)
	if err != nil {
		return fmt.Errorf("requesting storage counters: %v", err)
	}
	s.Gauges, err = storage.ListGauges(ctx)
	if err != nil {
		return fmt.Errorf("requesting storage gauges: %v", err)
	}
	return nil
}
