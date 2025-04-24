package routing

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/stepkareserva/obsermon/internal/server/metrics/handlers"
	"github.com/stepkareserva/obsermon/internal/server/middleware"
	"go.uber.org/zap"
)

type Routing struct {
	metrics *handlers.MetricsHandler
	log     *zap.Logger
}

func New(s Service, log *zap.Logger) (*Routing, error) {
	if s == nil {
		return nil, fmt.Errorf("service not exist")
	}
	metrics, err := handlers.New(s, log)
	if err != nil {
		return nil, fmt.Errorf("metrics handlers: %w", err)
	}
	return &Routing{
		metrics: metrics,
		log:     log,
	}, nil
}

func (r *Routing) Handler(ctx context.Context) (http.Handler, error) {
	if r == nil {
		return nil, fmt.Errorf("routing not exists")
	}

	router := chi.NewRouter()
	if r.log != nil {
		router.Use(middleware.Logger(r.log))
	}
	router.Use(middleware.Compression(r.log))
	router.Use(middleware.Buffering(r.log))

	metricsHandler, err := r.metrics.Handler(ctx)
	if err != nil {
		return nil, fmt.Errorf("metrics handler: %w", err)
	}
	router.Mount("/", metricsHandler)

	return router, nil
}
