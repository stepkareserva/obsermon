package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/stepkareserva/obsermon/cmd/server/internal/metrics/handlers"
	"github.com/stepkareserva/obsermon/cmd/server/internal/metrics/server"
	"github.com/stepkareserva/obsermon/cmd/server/internal/metrics/storage"
)

func main() {
	// initialize storage, controller, handler
	storage := storage.NewMemStorage()
	server, err := server.NewServer(storage)
	if err != nil {
		log.Fatal(err)
		return
	}
	updateHandler, err := handlers.UpdateHandler(server)
	if err != nil {
		log.Fatal(err)
		return
	}

	r := chi.NewRouter()
	r.Mount("/update", updateHandler)

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
