package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/stepkareserva/obsermon/cmd/agent/internal/client"
	"github.com/stepkareserva/obsermon/cmd/agent/internal/config"
	"github.com/stepkareserva/obsermon/cmd/agent/internal/watchdog"
)

func main() {
	// reading params
	cfg := config.ParseConfig()
	if err := cfg.Validate(); err != nil {
		log.Print(err)
		return
	}

	// metrics client
	metricsClient, err := client.NewMetricsClient(cfg.Endpoint)
	if err != nil {
		log.Printf("metrics client initialization error: %v", err)
		return
	}

	// watchdog
	watchdogParams := watchdog.WatchdogParams{
		PollInterval:        cfg.PollInterval,
		ReportInterval:      cfg.ReportInterval,
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
