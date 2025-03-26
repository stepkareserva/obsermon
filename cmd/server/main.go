package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/stepkareserva/obsermon/cmd/server/internal/config"
	"github.com/stepkareserva/obsermon/cmd/server/internal/metrics/handlers"
	"github.com/stepkareserva/obsermon/cmd/server/internal/metrics/server"
	"github.com/stepkareserva/obsermon/cmd/server/internal/metrics/storage"
)

func main() {
	// reading params
	cfg, err := readConfig()
	if err != nil {
		log.Println(err)
		return
	}

	// initialize storage and controller
	storage := storage.NewMemStorage()
	server, err := server.NewServer(storage)
	if err != nil {
		log.Fatal(err)
		return
	}

	// initialize handler
	handler, err := createHandler(server)
	if err != nil {
		log.Fatal(err)
		return
	}

	// run server
	log.Printf("Server is running on %s", cfg.Endpoint)
	log.Fatal(http.ListenAndServe(cfg.Endpoint, handler))
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

func createHandler(s *server.Server) (http.Handler, error) {

	updateHandler, err := handlers.UpdateHandler(s)
	if err != nil {
		return nil, err
	}
	valueHandler, err := handlers.ValueHandler(s)
	if err != nil {
		return nil, err
	}

	valuesHandler, err := handlers.ValuesHandler(s)
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()
	r.Mount("/update", updateHandler)
	r.Mount("/value", valueHandler)
	r.Mount("/", valuesHandler)

	return r, nil
}
