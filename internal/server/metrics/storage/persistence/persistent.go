package persistence

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/stepkareserva/obsermon/internal/models"
	"github.com/stepkareserva/obsermon/internal/server/metrics/service"
	"go.uber.org/zap"
)

type Config struct {
	StateStorage  StateStorage
	StoreInterval time.Duration
	Restore       bool
}

type Storage struct {
	service.Storage
	sstorage StateStorage

	saveCh chan time.Time
	cancel context.CancelFunc
	wg     sync.WaitGroup

	logger *zap.Logger
}

func New(cfg Config, base service.Storage, logger *zap.Logger) (*Storage, error) {
	if base == nil {
		return nil, fmt.Errorf("base service is nil")
	}
	if cfg.StateStorage == nil {
		return nil, fmt.Errorf("state io service is nil")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger is nil")
	}
	if cfg.Restore {
		state, err := cfg.StateStorage.LoadState()
		if err != nil {
			logger.Warn("state loading: %w", zap.Error(err))
		} else if err := state.Export(context.TODO(), base); err != nil {
			logger.Warn("state exporting: %w", zap.Error(err))
		}
	}

	ctx, cancel := context.WithCancel(context.TODO())

	// saveCh could be chan of struct{} but to unify storing process use chan of time,
	// because in some cases storing are based on ticker (which use chan of time)
	// and in other cases storing are based on our channel.
	// unoptimal? defenetly. simple? yes.
	storage := &Storage{
		Storage:  base,
		sstorage: cfg.StateStorage,
		saveCh:   make(chan time.Time),
		cancel:   cancel,
		logger:   logger,
	}

	// run storing loop, sync or async
	storage.wg.Add(1)
	go func() {
		defer storage.wg.Done()
		storage.runStoringLoop(ctx, cfg.StoreInterval)
	}()

	return storage, nil
}

func (s *Storage) Close() error {
	s.cancel()
	s.wg.Wait()
	return nil
}

func (s *Storage) SetGauge(ctx context.Context, val models.Gauge) error {
	if s == nil || s.Storage == nil {
		return fmt.Errorf("storage not exists")
	}
	if err := s.Storage.SetGauge(ctx, val); err != nil {
		return err
	}
	s.onModify()
	return nil
}

func (s *Storage) SetGauges(ctx context.Context, vals models.GaugesList) error {
	if s == nil || s.Storage == nil {
		return fmt.Errorf("storage not exists")
	}
	if err := s.Storage.SetGauges(ctx, vals); err != nil {
		return err
	}
	s.onModify()
	return nil
}

func (s *Storage) ReplaceGauges(ctx context.Context, val models.GaugesList) error {
	if s == nil || s.Storage == nil {
		return fmt.Errorf("storage not exists")
	}
	if err := s.Storage.ReplaceGauges(ctx, val); err != nil {
		return err
	}
	s.onModify()
	return nil
}

func (s *Storage) UpdateCounter(ctx context.Context, val models.Counter) (*models.Counter, error) {
	if s == nil || s.Storage == nil {
		return nil, fmt.Errorf("storage not exists")
	}
	updated, err := s.Storage.UpdateCounter(ctx, val)
	if err != nil {
		return nil, err
	}
	s.onModify()
	return updated, nil
}

func (s *Storage) UpdateCounters(ctx context.Context, vals models.CountersList) (models.CountersList, error) {
	if s == nil || s.Storage == nil {
		return nil, fmt.Errorf("storage not exists")
	}
	updated, err := s.Storage.UpdateCounters(ctx, vals)
	if err != nil {
		return nil, err
	}
	s.onModify()
	return updated, nil
}

func (s *Storage) ReplaceCounters(ctx context.Context, val models.CountersList) error {
	if s == nil || s.Storage == nil {
		return fmt.Errorf("storage not exists")
	}
	if err := s.Storage.ReplaceCounters(ctx, val); err != nil {
		return err
	}
	s.onModify()
	return nil
}

func (s *Storage) runStoringLoop(ctx context.Context, interval time.Duration) {
	if interval > 0 {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		s.channelStoringLoop(ctx, ticker.C)
	} else {
		s.channelStoringLoop(ctx, s.saveCh)
	}
}

func (s *Storage) channelStoringLoop(ctx context.Context, ch <-chan time.Time) {
	for {
		select {
		case <-ch:
			if err := s.storeState(ctx); err != nil {
				s.logger.Error("store state", zap.Error(err))
			}
		case <-ctx.Done():
			if err := s.storeState(ctx); err != nil {
				s.logger.Error("store state", zap.Error(err))
			}
			return
		}
	}
}

func (s *Storage) storeState(ctx context.Context) error {
	var state State
	if err := state.Import(ctx, s.Storage); err != nil {
		return fmt.Errorf("storage state request: %w", err)
	}
	if err := s.sstorage.StoreState(state); err != nil {
		return fmt.Errorf("storage state storing: %w", err)
	}
	return nil
}

func (s *Storage) onModify() {
	// write to save channel current time
	select {
	case s.saveCh <- time.Now():
	default:
	}
}
