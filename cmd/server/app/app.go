package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/stepkareserva/obsermon/internal/server/config"
	"github.com/stepkareserva/obsermon/internal/server/metrics/handlers"
	"github.com/stepkareserva/obsermon/internal/server/metrics/persistence"
	"github.com/stepkareserva/obsermon/internal/server/metrics/service"
	"github.com/stepkareserva/obsermon/internal/server/metrics/storage"
	"github.com/stepkareserva/obsermon/internal/server/server"
	"go.uber.org/zap"
)

type App struct {
	cfg     *config.Config
	log     *zap.Logger
	service *persistence.Service
	server  *server.Server
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewApp(cfg *config.Config, log *zap.Logger) (*App, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config not exists")
	}
	if log == nil {
		log = zap.NewNop()
	}
	service, err := initService(cfg, log)
	if err != nil {
		return nil, fmt.Errorf("service init: %w", err)
	}
	server, err := initServer(cfg, service, log)
	if err != nil {
		return nil, fmt.Errorf("server init: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &App{
		cfg:     cfg,
		log:     log,
		service: service,
		server:  server,
		ctx:     ctx,
		cancel:  cancel,
	}, nil
}

func (a *App) Run() (err error) {
	if a == nil {
		return fmt.Errorf("app not exists")
	}

	// start persistent service
	a.service.Start()
	defer func() {
		a.service.Stop()
		a.log.Info("service stoppeed")
	}()

	// start server
	if err := a.server.Start(); err != nil {
		return fmt.Errorf("server starting: %w", err)
	}
	defer func() {
		if stopErr := a.server.Stop(); stopErr != nil {
			err = errors.Join(err, stopErr)
		} else {
			a.log.Info("server stopped")
		}
	}()

	a.log.Info("server is running",
		zap.String("endpoint", a.cfg.Endpoint),
		zap.String("storage", a.cfg.FileStoragePath),
	)

	<-a.ctx.Done()

	a.log.Info("shutting down gracefully...")

	return nil
}

func (a *App) Stop() {
	a.cancel()
}

func initService(cfg *config.Config, log *zap.Logger) (*persistence.Service, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config not exists")
	}
	if log == nil {
		return nil, fmt.Errorf("log not exists")
	}

	// storage and service
	storage := storage.NewMemStorage()
	service, err := service.New(storage)
	if err != nil {
		return nil, fmt.Errorf("service creation: %w", err)
	}

	// wrap onto persistent object
	stateStorage := persistence.NewJSONStateStorage(cfg.FileStoragePath)
	persistenceCfg := persistence.ServiceConfig{
		Base:          service,
		StateStorage:  &stateStorage,
		Restore:       cfg.Restore,
		StoreInterval: cfg.StoreInterval(),
		Logger:        log,
	}
	persistenceService, err := persistence.New(persistenceCfg)
	if err != nil {
		return nil, fmt.Errorf("persistent service: %w", err)
	}

	return persistenceService, nil
}

func initServer(cfg *config.Config, service handlers.Service, log *zap.Logger) (*server.Server, error) {
	server, err := server.New(cfg, service, log)
	if err != nil {
		return nil, fmt.Errorf("server init: %w", err)
	}
	return server, nil
}
