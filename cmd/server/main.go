package main

import (
	"context"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/stepkareserva/obsermon/internal/server/config"
	"github.com/stepkareserva/obsermon/internal/server/logging"
	"github.com/stepkareserva/obsermon/internal/server/metrics/handlers"
	"github.com/stepkareserva/obsermon/internal/server/metrics/persistence"
	"github.com/stepkareserva/obsermon/internal/server/metrics/service"
	"github.com/stepkareserva/obsermon/internal/server/metrics/storage"
	"go.uber.org/zap"
)

func main() {
	// create log. use std log to log log errors,
	// because who log the log
	log, err := logging.New(logging.LevelDev)
	if err != nil {
		stdlog.Print(err)
		return
	}
	defer func() {
		if err := log.Sync(); err != nil {
			stdlog.Print(err)
		}
	}()

	// create gentle cancelling to context
	ctx, err := gracefulCancellingCtx(log)
	if err != nil {
		log.Error("graceful cancelling init", zap.Error(err))
	}

	// load and validate config
	cfg, err := loadConfig()
	if err != nil {
		log.Error("config loading", zap.Error(err))
		return
	}

	// initialize service
	service, err := initService(cfg, log)
	if err != nil {
		log.Error("service initialization", zap.Error(err))
		return
	}

	// run server in goroutine
	server, err := runServer(ctx, service, cfg, log)
	if err != nil {
		log.Error("server starting", zap.Error(err))
		return
	}

	// wait for cancel
	<-ctx.Done()

	// shutdown server
	if err = shutdown(server, service, log); err != nil {
		log.Error("shutdown", zap.Error(err))
	}
}

func gracefulCancellingCtx(log *zap.Logger) (context.Context, error) {
	if log == nil {
		return nil, fmt.Errorf("log not exists")
	}

	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		log.Info("Press Ctrl+C to stop the agent...")
		sig := <-sigChan
		log.Info(fmt.Sprintf("interruption signal received: %v, shutting down server...", sig))
		cancel()
	}()
	return ctx, nil
}

func loadConfig() (*config.Config, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("config loading: %v", err)
	}
	if err = config.Validate(*cfg); err != nil {
		return nil, fmt.Errorf("config validation: %v", err)
	}
	return cfg, nil
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

func runServer(
	ctx context.Context,
	service handlers.Service,
	cfg *config.Config,
	log *zap.Logger,
) (*http.Server, error) {
	if service == nil {
		return nil, fmt.Errorf("service not exists")
	}
	if log == nil {
		return nil, fmt.Errorf("log not exists")
	}
	if cfg == nil {
		return nil, fmt.Errorf("config not exists")
	}

	// create handlers
	handler, err := handlers.New(ctx, service, log)
	if err != nil {
		return nil, fmt.Errorf("handlers initialization: %w", err)
	}

	// run server in goroutine
	server := http.Server{Addr: cfg.Endpoint, Handler: handler}
	log.Info("server is running",
		zap.String("endpoint", cfg.Endpoint),
		zap.String("storage", cfg.FileStoragePath),
	)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Error("server listening", zap.Error(err))
		}
	}()

	return &server, nil
}

func shutdown(server *http.Server, service *persistence.Service, log *zap.Logger) error {
	if server == nil {
		return fmt.Errorf("server not exists")
	}
	if service == nil {
		return fmt.Errorf("service not exists")
	}
	if log == nil {
		return fmt.Errorf("log not exists")
	}

	// cancel server
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(context); err != nil {
		log.Error("server shutdown", zap.Error(err))
	} else {
		log.Info("server stopped")
	}

	// cancel service
	if err := service.Close(); err != nil {
		log.Error("service stopping", zap.Error(err))
	} else {
		log.Info("service stopped")
	}

	return nil
}
