package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/stepkareserva/obsermon/internal/server/config"
	"github.com/stepkareserva/obsermon/internal/server/logging"
	"github.com/stepkareserva/obsermon/internal/server/metrics/handlers"
	"github.com/stepkareserva/obsermon/internal/server/metrics/memstorage"
	"github.com/stepkareserva/obsermon/internal/server/metrics/persistence"
	"github.com/stepkareserva/obsermon/internal/server/metrics/service"
	"go.uber.org/zap"
)

func main() {
	// load and validate config
	cfg, err := loadConfig()
	if err != nil {
		stdlog.Printf("config loading: %v", err)
		return
	}

	// create log. use std log to log log errors,
	// because who log the log
	log, err := logging.New(cfg.Mode)
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

	// initialize storage
	storage, err := initStorage(cfg, log)
	if err != nil {
		log.Error("storage initialization", zap.Error(err))
	}

	// initialize service
	service, err := initService(cfg, storage, log)
	if err != nil {
		shutdown(nil, storage, log)
		log.Error("service initialization", zap.Error(err))
		return
	}

	// run server in goroutine
	server, err := runServer(ctx, service, cfg, log)
	if err != nil {
		shutdown(server, storage, log)
		log.Error("server starting", zap.Error(err))
		return
	}

	// wait for cancel
	<-ctx.Done()

	// shutdown server
	shutdown(server, storage, log)
}

func gracefulCancellingCtx(log *zap.Logger) (context.Context, error) {
	if log == nil {
		return nil, fmt.Errorf("log not exists")
	}

	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		log.Info("Press Ctrl+C to stop the server...")
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

func initStorage(cfg *config.Config, log *zap.Logger) (service.Storage, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config not exists")
	}
	if log == nil {
		return nil, fmt.Errorf("log not exists")
	}

	// storage
	storage := memstorage.New()

	// wrap onto persistent storage
	stateStorage := persistence.NewJSONStateStorage(cfg.FileStoragePath)
	persistenceCfg := persistence.Config{
		StateStorage:  &stateStorage,
		Restore:       cfg.Restore,
		StoreInterval: cfg.StoreInterval(),
	}
	persistentStorage, err := persistence.New(persistenceCfg, storage, log)
	if err != nil {
		return nil, fmt.Errorf("persistent service: %w", err)
	}

	return persistentStorage, nil
}

func initService(cfg *config.Config, storage service.Storage, log *zap.Logger) (*service.Service, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config not exists")
	}
	if storage == nil {
		return nil, fmt.Errorf("storage not exists")
	}
	if log == nil {
		return nil, fmt.Errorf("log not exists")
	}

	// service
	service, err := service.New(storage)
	if err != nil {
		return nil, fmt.Errorf("service creation: %w", err)
	}

	return service, nil
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

func shutdown(server *http.Server, storage service.Storage, log *zap.Logger) {
	if log == nil {
		log.Error("log not exists")
		return
	}

	// cancel server
	if server != nil {
		context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(context); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server shutdown", zap.Error(err))
		} else {
			log.Info("server stopped")
		}
	}

	// cancel storage, if it can be cancelled
	if storage != nil {
		if c, ok := storage.(io.Closer); ok {
			if err := c.Close(); err != nil {
				log.Error("storage closing", zap.Error(err))
			} else {
				log.Info("storage closed")
			}
		} else {
			log.Info("storage does not implement io.Closer")
		}
	}
}
