package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/stepkareserva/obsermon/cmd/agent/internal/client"
	"github.com/stepkareserva/obsermon/cmd/agent/internal/watchdog"
)

func main() {
	// params
	pollInterval := 2 * time.Second
	updateInterval := 10 * time.Second
	metricsServerURL := "http://localhost:8080"

	// metrics client
	metricsClient, err := client.NewMetricsClient(metricsServerURL)
	if err != nil {
		log.Printf("metrics client initialization error: %v", err)
		return
	}

	// watchdog
	watchdogParams := watchdog.WatchdogParams{
		PollInterval:        pollInterval,
		UpdateInterval:      updateInterval,
		MetricsServerClient: metricsClient,
	}
	watchdog := watchdog.NewWatchdog(watchdogParams)

	// run watchdog in goroutine
	go func() {
		watchdog.Start()
		log.Println("Running watchdog goroutine...")
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
