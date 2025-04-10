package main

import (
	"log"
	"net/http"

	"github.com/stepkareserva/obsermon/internal/server/config"
	"github.com/stepkareserva/obsermon/internal/server/metrics/handlers"
	"github.com/stepkareserva/obsermon/internal/server/metrics/service"
	"github.com/stepkareserva/obsermon/internal/server/metrics/storage"
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

	// initialize storage and controller
	storage := storage.NewMemStorage()
	service, err := service.New(storage)
	if err != nil {
		log.Printf("service initialization: %v", err)
		return
	}

	// initialize handler
	handler, err := handlers.New(service)
	if err != nil {
		log.Printf("handlers initialization: %v", err)
		return
	}

	// run server
	log.Printf("Server is running on %s", cfg.Endpoint)
	err = http.ListenAndServe(cfg.Endpoint, handler)
	if err != nil {
		log.Printf("server running: %v", err)
		return
	}
}
