package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/stepkareserva/obsermon/internal/server/config"
	"github.com/stepkareserva/obsermon/internal/server/handlers"
	"github.com/stepkareserva/obsermon/internal/server/metrics/persistence"
	"github.com/stepkareserva/obsermon/internal/server/metrics/service"
	"github.com/stepkareserva/obsermon/internal/server/metrics/storage"
	"github.com/stepkareserva/obsermon/internal/server/server"
	"go.uber.org/zap"
)

type App struct {
	cfg *config.Config
	log *zap.Logger

	storage service.Storage
	service *service.Service
	server  *server.Server

	ctx    context.Context
	cancel context.CancelFunc
}

func NewApp(cfg *config.Config, log *zap.Logger) (*App, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config not exists")
	}
	if log == nil {
		log = zap.NewNop()
	}
	storage, err := initStorage(cfg, log)
	if err != nil {
		return nil, fmt.Errorf("storage init: %w", err)
	}
	service, err := initService(cfg, storage, log)
	if err != nil {
		return nil, fmt.Errorf("service init: %w", err)
	}
	server, err := initServer(cfg, service, log)
	if err != nil {
		return nil, fmt.Errorf("server init: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &App{
		cfg: cfg,
		log: log,

		storage: storage,
		service: service,
		server:  server,

		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (a *App) Shutdown() {
	a.cancel()
}

func (a *App) Run() (err error) {
	if a == nil {
		return fmt.Errorf("app not exists")
	}

	// close storage
	defer func() {
		if stopErr := a.storage.Close(); err != nil {
			err = errors.Join(err, stopErr)
		} else {
			a.log.Info("storage stopped")
		}
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

func initStorage(cfg *config.Config, log *zap.Logger) (service.Storage, error) {
	// memory storage
	storage := storage.NewMemStorage()

	// add persistence
	stateStorage := persistence.NewJSONStateStorage(cfg.FileStoragePath)
	persistenceCfg := persistence.StorageConfig{
		Base:          storage,
		StateStorage:  &stateStorage,
		Restore:       cfg.Restore,
		StoreInterval: cfg.StoreInterval(),
		Logger:        log,
	}
	persistentStorage, err := persistence.New(persistenceCfg)
	if err != nil {
		return nil, fmt.Errorf("create persistent storage: %w", err)
	}

	return persistentStorage, nil
}

func initService(cfg *config.Config, storage service.Storage, log *zap.Logger) (*service.Service, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config not exists")
	}
	if log == nil {
		return nil, fmt.Errorf("log not exists")
	}
	if storage == nil {
		return nil, fmt.Errorf("storage not exists")
	}

	service, err := service.New(storage)
	if err != nil {
		return nil, fmt.Errorf("service creation: %w", err)
	}

	return service, nil
}

func initServer(cfg *config.Config, service handlers.Service, log *zap.Logger) (*server.Server, error) {
	server, err := server.New(cfg, service, log)
	if err != nil {
		return nil, fmt.Errorf("server init: %w", err)
	}
	return server, nil
}
