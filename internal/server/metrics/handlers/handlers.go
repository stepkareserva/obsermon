package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func New(s Service) (http.Handler, error) {
	updateHandler, err := UpdateHandler(s)
	if err != nil {
		return nil, err
	}
	valueHandler, err := ValueHandler(s)
	if err != nil {
		return nil, err
	}

	valuesHandler, err := ValuesHandler(s)
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()
	r.Mount("/update", updateHandler)
	r.Mount("/value", valueHandler)
	r.Mount("/", valuesHandler)

	return r, nil
}
