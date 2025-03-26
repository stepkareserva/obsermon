package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/stepkareserva/obsermon/cmd/agent/internal/client"
	"github.com/stepkareserva/obsermon/cmd/agent/internal/config"
	"github.com/stepkareserva/obsermon/cmd/agent/internal/watchdog"
)

func main() {
	// reading params from command line
	cfg, err := readConfig()
	if err != nil {
		log.Println(err)
		return
	}

	// metrics client
	metricsClient, err := client.NewMetricsClient(cfg.EndpointURL())
	if err != nil {
		log.Printf("metrics client initialization error: %v", err)
		return
	}

	// watchdog
	watchdogParams := watchdog.WatchdogParams{
		PollInterval:        time.Duration(cfg.PollInterval()),
		ReportInterval:      time.Duration(cfg.ReportInterval()),
		MetricsServerClient: metricsClient,
	}
	watchdog := watchdog.NewWatchdog(watchdogParams)

	// run watchdog in goroutine
	go func() {
		log.Println("Running watchdog goroutine...")
		watchdog.Start()
	}()

	// wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Press Ctrl+C to stop the agent...")
	<-sigChan

	// shut down the server
	log.Println("Received interrupt signal. Shutting down agent...")
	watchdog.Stop()
	log.Println("Agent shut down")
}

func readConfig() (*config.Config, error) {
	var cfg config.Config
	if err := cfg.ParseCommandLine(); err != nil {
		return nil, fmt.Errorf("error parsing command line: %w", err)
	}
	if err := cfg.ParseEnv(); err != nil {
		return nil, fmt.Errorf("error parsing env: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return &cfg, nil
}
