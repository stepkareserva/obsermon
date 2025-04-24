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
	sstorage      StateStorage
	logger        *zap.Logger
	storeInterval time.Duration

	saveCh chan struct{}
	stopCh chan struct{}
	wg     sync.WaitGroup
}

func New(cfg ServiceConfig) (*Service, error) {
	if err := checkCreationParams(cfg); err != nil {
		return nil, err
	}
	service := &Service{
		BaseService:   cfg.Base,
		sstorage:      cfg.StateStorage,
		logger:        cfg.Logger,
		storeInterval: cfg.StoreInterval,
		saveCh:        make(chan struct{}),
		stopCh:        make(chan struct{}),
	}

	// restore if required and possible
	if cfg.Restore {
		if err := service.restore(); err != nil {
			// fail but ok ignore it
			service.logger.Warn("service config restoring", zap.Error(err))
		}
	}

	return service, nil
}

func (s *Service) Start() {
	// run storing loop, sync or async
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.runStoringLoop(s.storeInterval)
	}()
}

func (s *Service) Stop() {
	if s == nil {
		return
	}

	select {
	case s.stopCh <- struct{}{}:
	default:
		s.wg.Wait()
	}
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

func (s *Service) store() {
	state, err := s.BaseService.GetState()
	if err != nil {
		s.logger.Error("service state request", zap.Error(err))
		return
	}
	if err = s.sstorage.StoreState(*state); err != nil {
		s.logger.Error("service state storing", zap.Error(err))
	}
	s.logger.Info("service state stored")
}

func (s *Service) UpdateGauge(val models.Gauge) (*models.Gauge, error) {
	updated, err := s.BaseService.UpdateGauge(val)

	if err == nil {
		s.onUpdate()
	}

	return updated, err
}

func (s *Service) UpdateCounter(val models.Counter) (*models.Counter, error) {
	updated, err := s.BaseService.UpdateCounter(val)

	if err == nil {
		s.onUpdate()
	}

	return updated, err
}

func (s *Service) UpdateMetric(val models.Metrics) (*models.Metrics, error) {
	updated, err := s.BaseService.UpdateMetric(val)

	if err == nil {
		s.onUpdate()
	}

	return updated, err
}

func (s *Service) onUpdate() {
	select {
	case s.saveCh <- struct{}{}:
	default:
	}
}
