package main

import (
	"context"
	stdlog "log"
	"net/http"

	"github.com/stepkareserva/obsermon/internal/server/config"
	"github.com/stepkareserva/obsermon/internal/server/logging"
	"github.com/stepkareserva/obsermon/internal/server/metrics/handlers"
	"github.com/stepkareserva/obsermon/internal/server/metrics/service"
	"github.com/stepkareserva/obsermon/internal/server/metrics/storage"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	// create logger
	logger, err := logging.NewZapLogger(logging.LevelProd)
	if err != nil {
		stdlog.Print(err)
		return
	}
	defer logger.Sync()

	// add logger to context
	ctx = logging.WithLogger(ctx, logger)

	// load and validate config
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("config loading", zap.Error(err))
		return
	}
	if err = config.Validate(*cfg); err != nil {
		logger.Error("config validation", zap.Error(err))
		return
	}

	// initialize storage and controller
	storage := storage.NewMemStorage()
	service, err := service.New(storage)
	if err != nil {
		logger.Error("service initialization", zap.Error(err))
		return
	}

	// initialize handler
	handler, err := handlers.New(ctx, service)
	if err != nil {
		logger.Error("handlers initialization", zap.Error(err))
		return
	}

	// run server
	logger.Info("server is running",
		zap.String("endpoint", cfg.Endpoint),
		zap.String("storage", cfg.FileStoragePath),
	)
	err = http.ListenAndServe(cfg.Endpoint, handler)
	if err != nil {
		logger.Info("server starting", zap.Error(err))
		return
	}
}
