package routing

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	dbhandlers "github.com/stepkareserva/obsermon/internal/server/db/handlers"
	"github.com/stepkareserva/obsermon/internal/server/interfaces/database"
	mhandlers "github.com/stepkareserva/obsermon/internal/server/metrics/handlers"
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

func (r *Routing) AddMetricsHandlers(ctx context.Context, s mhandlers.Service) error {
	if r == nil || r.router == nil {
		return fmt.Errorf("routing not exists")
	}

	metricsHandlers, err := mhandlers.New(s, r.log)
	if err != nil {
		return fmt.Errorf("metrics handlers: %w", err)
	}

	if err := metricsHandlers.RegisterRoutes(ctx, r.router); err != nil {
		return fmt.Errorf("register metrics routes: %w", err)
	}

	return nil
}

func (r *Routing) AddDatabaseHandlers(ctx context.Context, db database.Database) error {
	if r == nil || r.router == nil {
		return fmt.Errorf("routing not exists")
	}

	dbHandlers, err := dbhandlers.New(db, r.log)
	if err != nil {
		return fmt.Errorf("metrics handlers: %w", err)
	}

	if err := dbHandlers.RegisterRoutes(ctx, r.router); err != nil {
		return fmt.Errorf("register database routes: %w", err)
	}

	return nil
}
