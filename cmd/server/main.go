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

	app, err := app.New(context.TODO(), *cfg, log)
	if err != nil {
		log.Error("app init", zap.Error(err))
	}
	defer func() {
		if err := app.Close(); err != nil {
			log.Error("app closing", zap.Error(err))
		}
	}()

	// create gentle cancelling to context
	ctx, err := gracefulCancellingCtx(log)
	if err != nil {
		log.Error("graceful cancelling init", zap.Error(err))
	}

	if err := app.Run(ctx); err != nil {
		log.Error("app running", zap.Error(err))
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
