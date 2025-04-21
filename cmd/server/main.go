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
	// create logger
	logger, err := logging.NewZapLogger(logging.LevelDev)
	if err != nil {
		stdlog.Print(err)
		return
	}
	defer logger.Sync()

	// create gentle cancelling to context
	ctx := gracefulCancellingCtx(logger)
	// add logger to context
	ctx = logging.WithLogger(ctx, logger)

	// load and validate config
	cfg, err := loadConfig()
	if err != nil {
		logger.Error("config loading", zap.Error(err))
		return
	}

	// initialize service
	service, err := initService(cfg, logger)
	if err != nil {
		logger.Error("service initialization", zap.Error(err))
		return
	}

	// run server in goroutine
	server, err := runServer(service, *cfg, ctx)
	if err != nil {
		logger.Error("server starting", zap.Error(err))
		return
	}

	// wait for cancel
	<-ctx.Done()

	// shutdown server
	shutdown(server, service, logger)
}

func gracefulCancellingCtx(logger *zap.Logger) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		logger.Info("Press Ctrl+C to stop the agent...")
		sig := <-sigChan
		logger.Info(fmt.Sprintf("interruption signal received: %v, shutting down server...", sig))
		cancel()
	}()
	return ctx
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

func initService(cfg *config.Config, logger *zap.Logger) (*persistence.Service, error) {
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
		Logger:        logger,
	}
	persistenceService, err := persistence.New(persistenceCfg)
	if err != nil {
		return nil, fmt.Errorf("persistent service: %w", err)
	}

	return persistenceService, nil
}

func runServer(service handlers.Service, cfg config.Config, ctx context.Context) (*http.Server, error) {
	logger := logging.FromContext(ctx)

	// create handlers
	handler, err := handlers.New(ctx, service)
	if err != nil {
		return nil, fmt.Errorf("handlers initialization: %w", err)
	}

	// run server in goroutine
	server := http.Server{Addr: cfg.Endpoint, Handler: handler}
	logger.Info("server is running",
		zap.String("endpoint", cfg.Endpoint),
		zap.String("storage", cfg.FileStoragePath),
	)

	go func() {
		server.ListenAndServe()
	}()

	return &server, nil
}

func shutdown(server *http.Server, service *persistence.Service, logger *zap.Logger) {
	// cancel server
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(context); err != nil {
		logger.Error("server shutdown", zap.Error(err))
	} else {
		logger.Info("server stopped")
	}
	// cancel service
	if err := service.Close(); err != nil {
		logger.Error("service stopping", zap.Error(err))
	} else {
		logger.Info("service stopped")
	}
}
