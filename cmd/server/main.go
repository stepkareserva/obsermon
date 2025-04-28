package main

import (
	"context"
	"fmt"
	stdlog "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/stepkareserva/obsermon/cmd/server/app"
	"github.com/stepkareserva/obsermon/internal/server/config"
	"github.com/stepkareserva/obsermon/internal/server/logging"
	"go.uber.org/zap"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		stdlog.Printf("config loading: %v", err)
		return
	}

	logger, err := logging.New(cfg.Mode)
	if err != nil {
		stdlog.Printf("logger creation: %v", err)
		return
	}

	app, err := app.NewApp(cfg, logger)
	if err != nil {
		logger.Error("app initialization", zap.Error(err))
		return
	}

	ctx, err := gracefulCancellingCtx(logger)
	if err != nil {
		logger.Error("graceful cancelling init", zap.Error(err))
	}

	go func() {
		<-ctx.Done()
		app.Shutdown()
	}()

	if err := app.Run(); err != nil {
		logger.Error("app running", zap.Error(err))
	}
	defer app.Shutdown()
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
