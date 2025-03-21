package main

import (
	"fmt"
	"log"
	"net/http"

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
	updateHandler, err := handlers.NewUpdateHandler(server)
	if err != nil {
		log.Fatal(err)
		return
	}

	mux := http.NewServeMux()

	mux.Handle("/update/gauge/", http.StripPrefix("/update/gauge/",
		http.HandlerFunc(updateHandler.UpdateGaugeHandler)))
	mux.Handle("/update/counter/", http.StripPrefix("/update/counter/",
		http.HandlerFunc(updateHandler.UpdateCounterHandler)))
	mux.Handle("/update/", http.StripPrefix("/update/",
		http.HandlerFunc(updateHandler.UpdateHandler)))

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
