package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/stepkareserva/obsermon/internal/server/middleware"
	"go.uber.org/zap"
)

func New(s Service, l *zap.Logger) (http.Handler, error) {
	if s == nil {
		return nil, fmt.Errorf("service not exist")
	}

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

	if l != nil {
		r.Use(middleware.Logger(l))
	}
	r.Use(middleware.Buffering())

	r.Mount("/update", updateHandler)
	r.Mount("/value", valueHandler)
	r.Mount("/", valuesHandler)

	return r, nil
}
