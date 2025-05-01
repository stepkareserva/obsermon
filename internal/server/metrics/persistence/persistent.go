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
	Restore       bool
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

	// restore if required and possible
	if cfg.Restore {
		if err := storage.loadState(); err != nil {
			// fail but ok ignore it
			logger.Warn("service config restoring", zap.Error(err))
		}
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
	if err = setStorageState(s.Storage, state); err != nil {
		return fmt.Errorf("storage state request: %w", err)
	}
	return nil
}

func (s *Storage) storeState() error {
	state, err := getStorageState(s.Storage)
	if err != nil {
		return fmt.Errorf("storage state request: %w", err)
	}
	if err = s.sstorage.StoreState(*state); err != nil {
		return fmt.Errorf("storage state storing: %w", err)
	}
	return nil
}

func getStorageState(storage service.Storage) (*State, error) {
	if storage == nil {
		return nil, fmt.Errorf("storage not exists")
	}
	counters, err := storage.ListCounters()
	if err != nil {
		return nil, fmt.Errorf("getting counters: %w", err)
	}
	gauges, err := storage.ListGauges()
	if err != nil {
		return nil, fmt.Errorf("getting gauges: %w", err)
	}
	return &State{
		Counters: counters,
		Gauges:   gauges,
	}, nil
}

func setStorageState(storage service.Storage, state *State) error {
	if storage == nil {
		return fmt.Errorf("storage not exists")
	}
	if state == nil {
		return fmt.Errorf("state not exists")
	}
	if err := storage.ReplaceCounters(state.Counters); err != nil {
		return fmt.Errorf("replacing storage counters: %w", err)
	}
	if err := storage.ReplaceGauges(state.Gauges); err != nil {
		return fmt.Errorf("replacing storage gauges: %w", err)
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
