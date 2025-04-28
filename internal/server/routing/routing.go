package routing

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/stepkareserva/obsermon/internal/server/handlers"
	"github.com/stepkareserva/obsermon/internal/server/middleware"
	"go.uber.org/zap"
)

type Routing struct {
	hMetrics  *handlers.MetricsHandler
	hDatabase *handlers.DBHandler
	log       *zap.Logger
}

func New(s Service, db Database, log *zap.Logger) (*Routing, error) {
	if s == nil {
		return nil, fmt.Errorf("service not exist")
	}

	hMetrics, err := handlers.NewMetricsHandler(s, log)
	if err != nil {
		return nil, fmt.Errorf("metrics handlers: %w", err)
	}

	hDatabase, err := handlers.NewDBHandler(db, log)
	if err != nil {
		return nil, fmt.Errorf("database handlers: %w", err)
	}

	return &Routing{
		hMetrics:  hMetrics,
		hDatabase: hDatabase,
		log:       log,
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

	metricsHandler, err := r.hMetrics.Handler(ctx)
	if err != nil {
		return nil, fmt.Errorf("metrics handler: %w", err)
	}
	router.Mount("/", metricsHandler)

	databaseHandler := r.hDatabase.PingHandler(ctx)
	router.Mount("/ping", databaseHandler)

	return router, nil
}
