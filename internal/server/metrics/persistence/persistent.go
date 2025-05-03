package persistence

import (
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
}

type Storage struct {
	service.Storage
	sstorage StateStorage

	saveCh chan time.Time
	stopCh chan struct{}
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

	// saveCh could be chan of struct{} but to unify storing process use chan of time,
	// because in some cases storing are based on ticker (which use chan of time)
	// and in other cases storing are based on our channel.
	// unoptimal? defenetly. simple? yes.
	storage := &Storage{
		Storage:  base,
		sstorage: cfg.StateStorage,
		saveCh:   make(chan time.Time),
		stopCh:   make(chan struct{}),
		logger:   logger,
	}

	// run storing loop, sync or async
	storage.wg.Add(1)
	go func() {
		defer storage.wg.Done()
		storage.runStoringLoop(cfg.StoreInterval)
	}()

	return storage, nil
}

func (s *Storage) Close() error {
	s.stopCh <- struct{}{}
	s.wg.Wait()
	return nil
}

func (s *Storage) SetGauge(val models.Gauge) error {
	if s == nil || s.Storage == nil {
		return fmt.Errorf("storage not exists")
	}
	if err := s.Storage.SetGauge(val); err != nil {
		return err
	}
	s.onModify()
	return nil
}

func (s *Storage) ReplaceGauges(val models.GaugesList) error {
	if s == nil || s.Storage == nil {
		return fmt.Errorf("storage not exists")
	}
	if err := s.Storage.ReplaceGauges(val); err != nil {
		return err
	}
	s.onModify()
	return nil
}

func (s *Storage) SetCounter(val models.Counter) error {
	if s == nil || s.Storage == nil {
		return fmt.Errorf("storage not exists")
	}
	if err := s.Storage.SetCounter(val); err != nil {
		return err
	}
	s.onModify()
	return nil
}

func (s *Storage) ReplaceCounters(val models.CountersList) error {
	if s == nil || s.Storage == nil {
		return fmt.Errorf("storage not exists")
	}
	if err := s.Storage.ReplaceCounters(val); err != nil {
		return err
	}
	s.onModify()
	return nil
}

func (s *Storage) runStoringLoop(interval time.Duration) {
	if interval > 0 {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		s.channelStoringLoop(ticker.C)
	} else {
		s.channelStoringLoop(s.saveCh)
	}
}

func (s *Storage) channelStoringLoop(ch <-chan time.Time) {
	for {
		select {
		case <-ch:
			if err := s.storeState(); err != nil {
				s.logger.Error("store state", zap.Error(err))
			}
		case <-s.stopCh:
			if err := s.storeState(); err != nil {
				s.logger.Error("store state", zap.Error(err))
			}
			return
		}
	}
}

func (s *Storage) loadState() error {
	state, err := s.sstorage.LoadState()
	if err != nil {
		return fmt.Errorf("storage state loading: %w", err)
	}
	if err = state.Export(s.Storage); err != nil {
		return fmt.Errorf("storage state request: %w", err)
	}
	return nil
}

func (s *Storage) storeState() error {
	var state State
	if err := state.Import(s.Storage); err != nil {
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
