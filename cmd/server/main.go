package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/stepkareserva/obsermon/cmd/server/internal/config"
	"github.com/stepkareserva/obsermon/cmd/server/internal/metrics/handlers"
	"github.com/stepkareserva/obsermon/cmd/server/internal/metrics/server"
	"github.com/stepkareserva/obsermon/cmd/server/internal/metrics/storage"
)

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

func main() {
	// reading params
	cfg := config.ParseConfig()
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
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
