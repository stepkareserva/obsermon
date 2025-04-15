package persistence

import (
	"fmt"
	"sync"
	"time"

	"github.com/stepkareserva/obsermon/internal/models"
	"go.uber.org/zap"
)

type ServiceConfig struct {
	Base          BaseService
	StateStorage  StateStorage
	Restore       bool
	StoreInterval time.Duration
	Logger        *zap.Logger
}

type Service struct {
	BaseService
	sstorage StateStorage
	logger   *zap.Logger

	saveCh chan struct{}
	stopCh chan struct{}
	wg     sync.WaitGroup
}

func New(cfg ServiceConfig) (*Service, error) {
	if err := checkCreationParams(cfg); err != nil {
		return nil, err
	}
	service := &Service{
		BaseService: cfg.Base,
		sstorage:    cfg.StateStorage,
		logger:      cfg.Logger,
		saveCh:      make(chan struct{}),
		stopCh:      make(chan struct{}),
	}

	// restore if required
	if cfg.Restore {
		if err := service.restore(); err != nil {
			return nil, err
		}
	}

	// run storing loop, sync or async
	service.wg.Add(1)
	go func() {
		defer service.wg.Done()
		service.runStoringLoop(cfg.StoreInterval)
	}()

	return service, nil
}

func checkCreationParams(cfg ServiceConfig) error {
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

func (s *Service) Close() error {
	s.stopCh <- struct{}{}
	s.wg.Wait()
	return nil
}

func (s *Service) restore() error {
	state, err := s.sstorage.LoadState()
	if err != nil {
		return err
	}
	if err := s.BaseService.LoadState(*state); err != nil {
		return err
	}
	return nil
}

func (s *Service) runStoringLoop(interval time.Duration) {
	if interval > 0 {
		s.asyncStoringLoop(interval)
	} else {
		s.syncStoringLoop()
	}
}

func (s *Service) asyncStoringLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.store()
		case <-s.stopCh:
			s.store()
			return
		}
	}
}

func (s *Service) syncStoringLoop() {
	for {
		select {
		case <-s.saveCh:
			s.store()
		case <-s.stopCh:
			s.store()
			return
		}
	}
}

func (s *Service) store() error {
	state, err := s.BaseService.GetState()
	if err != nil {
		return err
	}
	if err = s.sstorage.StoreState(*state); err != nil {
		s.logger.Error("service state storing", zap.Error(err))
		return err
	}
	s.logger.Info("service state stored")
	return nil
}

func (s *Service) UpdateGauge(val models.Gauge) error {
	if err := s.BaseService.UpdateGauge(val); err != nil {
		return err
	}

	s.onUpdate()

	return nil
}

func (s *Service) UpdateCounter(val models.Counter) error {
	s.logger.Info("update counter")
	if err := s.BaseService.UpdateCounter(val); err != nil {
		return err
	}

	s.onUpdate()

	return nil
}

func (s *Service) UpdateMetric(val models.Metrics) error {
	if err := s.BaseService.UpdateMetric(val); err != nil {
		return err
	}

	s.onUpdate()

	return nil
}

func (s *Service) onUpdate() {
	select {
	case s.saveCh <- struct{}{}:
	default:
	}
}
