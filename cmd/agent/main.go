package main

import (
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
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
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
	watchdog, err := watchdog.NewWatchdog(watchdogParams)
	if err != nil {
		log.Printf("watchdog initialization error: %v", err)
		return
	}

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
