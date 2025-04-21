package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/stepkareserva/obsermon/internal/server/middleware"
	"go.uber.org/zap"
)

func New(ctx context.Context, s Service, log *zap.Logger) (http.Handler, error) {
	if s == nil {
		return nil, fmt.Errorf("service not exist")
	}

	updateHandler, err := UpdateHandler(s, log)
	if err != nil {
		return nil, err
	}
	valueHandler, err := ValueHandler(s, log)
	if err != nil {
		return nil, err
	}

	valuesHandler, err := ValuesHandler(s, log)
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()
	if log != nil {
		r.Use(middleware.Logger(log))
	}
	r.Use(middleware.Compression(log))
	r.Use(middleware.Buffering())

	r.Mount("/update", updateHandler)
	r.Mount("/value", valueHandler)
	r.Mount("/", valuesHandler)

	return r, nil
}
