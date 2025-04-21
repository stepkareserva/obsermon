package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/stepkareserva/obsermon/internal/server/logging"
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
	logger := logging.FromContext(ctx)
	if logger != nil {
		r.Use(middleware.Logger(logger))
	}
	r.Use(middleware.Compression(logger))
	r.Use(middleware.Buffering())

	r.Mount("/update", updateHandler)
	r.Mount("/value", valueHandler)
	r.Mount("/", valuesHandler)

	return r, nil
}
