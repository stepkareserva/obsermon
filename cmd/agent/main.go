package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/stepkareserva/obsermon/internal/agent/client"
	"github.com/stepkareserva/obsermon/internal/agent/config"
	"github.com/stepkareserva/obsermon/internal/agent/watchdog"
)

func main() {
	// load and validate config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("config loading: %v", err)
		return
	}
	if err = config.Validate(*cfg); err != nil {
		log.Printf("config validation: %v", err)
		return
	}

	// metrics client
	metricsClient, err := client.New(cfg.EndpointURL(), cfg.ReportSignKey)
	if err != nil {
		log.Printf("metrics client initialization: %v", err)
		return
	}

	// context to stop on interription
	ctx, cancel := context.WithCancel(context.Background())

	// goroutine to handle interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		log.Printf("Press Ctrl+C to stop the agent...")
		sig := <-sigChan
		log.Printf("interruption signal received: %v, shutting down agent...", sig)
		cancel()
	}()

	// watchdog
	watchdogParams := watchdog.WatchdogParams{
		PollInterval:        time.Duration(cfg.PollInterval()),
		ReportInterval:      time.Duration(cfg.ReportInterval()),
		MetricsServerClient: metricsClient,
	}
	watchdog, err := watchdog.New(watchdogParams)
	if err != nil {
		log.Printf("watchdog initialization: %v", err)
		return
	}
	watchdog.Start(ctx)

	log.Println("Agent shut down")
}
