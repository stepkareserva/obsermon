package persistence

import (
	"fmt"
	"sync"
	"time"

	"github.com/stepkareserva/obsermon/internal/models"
	"github.com/stepkareserva/obsermon/internal/server/metrics/service"
	"go.uber.org/zap"
)

type StorageConfig struct {
	Base          BaseStorage
	StateStorage  StateStorage
	Restore       bool
	StoreInterval time.Duration
	Logger        *zap.Logger
}

type Storage struct {
	service.Storage
	logger *zap.Logger

	saveCh chan struct{}
	stopCh chan struct{}
	wg     sync.WaitGroup
}

func New(cfg StorageConfig) (*Storage, error) {
	if err := checkCreationParams(cfg); err != nil {
		return nil, err
	}
	storage := &Storage{
		Storage: cfg.Base,
		logger:  cfg.Logger,
		saveCh:  make(chan struct{}),
		stopCh:  make(chan struct{}),
	}

	// restore if required and possible
	if cfg.Restore {
		if err := storage.load(cfg.StateStorage); err != nil {
			// fail but ok ignore it
			storage.logger.Warn("storage restoring", zap.Error(err))
		}
	}

	// run storing loop, sync or async
	storage.wg.Add(1)
	go func() {
		defer storage.wg.Done()
		storage.runStoringLoop(cfg.StateStorage, cfg.StoreInterval)
	}()

	return storage, nil
}

func (s *Storage) Close() error {
	if s == nil {
		return nil
	}

	select {
	case s.stopCh <- struct{}{}:
	default:
		s.wg.Wait()
	}

	return nil
}

func checkCreationParams(cfg StorageConfig) error {
	if cfg.Base == nil {
		return fmt.Errorf("base service is nil")
	}
	if cfg.StateStorage == nil {
		return fmt.Errorf("state io service is nil")
	}
	if cfg.Logger == nil {
		return fmt.Errorf("logger is nil")
	}
	return nil
}

func (s *Storage) runStoringLoop(storage StateStorage, interval time.Duration) {
	if interval > 0 {
		s.asyncStoringLoop(storage, interval)
	} else {
		s.syncStoringLoop(storage)
	}
}

func (s *Storage) asyncStoringLoop(storage StateStorage, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := s.store(storage); err != nil {
				s.logger.Error("metrics storing", zap.Error(err))
			}
		case <-s.stopCh:
			if err := s.store(storage); err != nil {
				s.logger.Error("metrics storing", zap.Error(err))
			}
			return
		}
	}
}

func (s *Storage) syncStoringLoop(storage StateStorage) {
	for {
		select {
		case <-s.saveCh:
			if err := s.store(storage); err != nil {
				s.logger.Error("metrics storing", zap.Error(err))
			}
		case <-s.stopCh:
			if err := s.store(storage); err != nil {
				s.logger.Error("metrics storing", zap.Error(err))
			}
			return
		}
	}
}

func (s *Storage) load(storage StateStorage) error {
	state, err := storage.LoadState()
	if err != nil {
		return fmt.Errorf("storage state loading: %w", err)
	}
	if err = s.setState(*state); err != nil {
		return fmt.Errorf("storage state applying: %w", err)
	}
	return nil
}

func (s *Storage) store(storage StateStorage) error {
	state, err := s.getState()
	if err != nil {
		return fmt.Errorf("service state request: %w", err)
	}
	if err = storage.StoreState(*state); err != nil {
		return fmt.Errorf("service state storing: %w", err)
	}
	return nil
}

func (s *Storage) getState() (*State, error) {
	var state State
	var err error

	if state.Counters, err = s.Storage.ListCounters(); err != nil {
		return nil, err
	}
	if state.Gauges, err = s.Storage.ListGauges(); err != nil {
		return nil, err
	}

	return &state, nil
}

func (s *Storage) setState(state State) error {
	if err := s.Storage.ReplaceCounters(state.Counters); err != nil {
		return err
	}
	if err := s.Storage.ReplaceGauges(state.Gauges); err != nil {
		return err
	}
	return nil
}

func (s *Storage) SetGauge(val models.Gauge) error {
	if err := s.Storage.SetGauge(val); err != nil {
		return err
	}

	s.onModify()
	return nil
}

func (s *Storage) ReplaceGauges(val models.GaugesList) error {
	if err := s.Storage.ReplaceGauges(val); err != nil {
		return err
	}

	s.onModify()
	return nil
}

func (s *Storage) SetCounter(val models.Counter) error {
	if err := s.Storage.SetCounter(val); err != nil {
		return err
	}

	s.onModify()
	return nil
}

func (s *Storage) ReplaceCounters(val models.CountersList) error {
	if err := s.Storage.ReplaceCounters(val); err != nil {
		return err
	}

	s.onModify()
	return nil
}

func (s *Storage) onModify() {
	select {
	case s.saveCh <- struct{}{}:
	default:
	}
}
