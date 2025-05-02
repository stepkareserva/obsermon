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
	router *chi.Mux
	log    *zap.Logger
}

func New(log *zap.Logger) Routing {
	if log == nil {
		log = zap.NewNop()
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger(log))
	r.Use(middleware.Compression(log))
	r.Use(middleware.Buffering(log))

	return Routing{
		router: r,
		log:    log,
	}
}

func (r *Routing) Handler() (http.Handler, error) {
	if r == nil || r.router == nil {
		return nil, fmt.Errorf("routing not exists")
	}
	return r.router, nil
}

func (r *Routing) AddMetricsHandlers(ctx context.Context, s handlers.Service) error {
	if r == nil || r.router == nil {
		return fmt.Errorf("routing not exists")
	}
	if s == nil {
		return fmt.Errorf("service not exists")
	}

	metricsHandlers, err := handlers.New(s, r.log)
	if err != nil {
		return fmt.Errorf("metrics handlers: %w", err)
	}

	if err := metricsHandlers.RegisterRoutes(ctx, r.router); err != nil {
		return fmt.Errorf("register metrics routes: %w", err)
	}

	return nil
}
