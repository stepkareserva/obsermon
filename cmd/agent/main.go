package main

import (
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
	metricsClient, err := client.New(cfg.EndpointURL())
	if err != nil {
		log.Printf("metrics client initialization: %v", err)
		return
	}

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

	// shut down agent
	log.Println("Received interrupt signal. Shutting down agent...")
	watchdog.Stop()
	log.Println("Agent shut down")
}
