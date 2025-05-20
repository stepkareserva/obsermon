package router

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/stepkareserva/obsermon/internal/server/http/constants"
	"github.com/stepkareserva/obsermon/internal/server/http/handlers"
	"github.com/stepkareserva/obsermon/internal/server/http/middleware"

	"go.uber.org/zap"
)

type Routing struct {
	router *chi.Mux
	log    *zap.Logger
}

func New(log *zap.Logger, s handlers.Service) (http.Handler, error) {
	if log == nil {
		log = zap.NewNop()
	}

	// add middleware
	r := chi.NewRouter()
	r.Use(middleware.Logger(log))
	r.Use(middleware.Compression(log))
	r.Use(middleware.Buffering(log))

	// register routes
	if err := addUpdateHandlers(r, s, log); err != nil {
		return nil, fmt.Errorf("update handlers: %w", err)
	}
	if err := addValueHandlers(r, s, log); err != nil {
		return nil, fmt.Errorf("value handlers: %w", err)
	}
	if err := addValuesHandlers(r, s, log); err != nil {
		return nil, fmt.Errorf("values handlers: %w", err)
	}
	if err := addPingHandlers(r, s, log); err != nil {
		return nil, fmt.Errorf("ping handlers: %w", err)
	}

	return r, nil
}

func addUpdateHandlers(r chi.Router, s handlers.Service, log *zap.Logger) error {
	updHandler, err := handlers.NewUpdateHandler(s, log)
	if err != nil {
		return fmt.Errorf("update handler creation: %w", err)
	}
	r.Route("/update", func(r chi.Router) {
		r.Post(fmt.Sprintf("/%s/{%s}/{%s}", constants.MetricGauge, constants.ChiName, constants.ChiValue),
			updHandler.UpdateGaugeURLHandler())
		r.Post(fmt.Sprintf("/%s/{%s}/{%s}", constants.MetricCounter, constants.ChiName, constants.ChiValue),
			updHandler.UpdateCounterURLHandler())
		r.Post(fmt.Sprintf("/{%s}/{%s}/{%s}", constants.ChiMetric, constants.ChiName, constants.ChiValue),
			updHandler.UpdateUnknownMetricURLHandler())
		r.Post("/",
			updHandler.UpdateMetricJSONHandler())
	})
	r.Route("/updates", func(r chi.Router) {
		r.Post("/",
			updHandler.UpdateMetricsJSONHandler())
	})

	return nil
}

func addValueHandlers(r chi.Router, s handlers.Service, log *zap.Logger) error {
	valHandler, err := handlers.NewValueHandler(s, log)
	if err != nil {
		return fmt.Errorf("value handler creation: %w", err)
	}
	r.Route("/value", func(r chi.Router) {
		r.Get(fmt.Sprintf("/%s/{%s}", constants.MetricGauge, constants.ChiName),
			valHandler.GaugeValueURLHandler())
		r.Get(fmt.Sprintf("/%s/{%s}", constants.MetricCounter, constants.ChiName),
			valHandler.CounterValueURLHandler())
		r.Get(fmt.Sprintf("/{%s}/{%s}", constants.ChiMetric, constants.ChiName),
			valHandler.UnknownMetricValueURLHandler())
		r.Post("/",
			valHandler.ValueMetricJSONHandler())
	})

	return nil
}

func addValuesHandlers(r chi.Router, s handlers.Service, log *zap.Logger) error {
	valsHandler, err := handlers.NewValuesHandler(s, log)
	if err != nil {
		return fmt.Errorf("values handler creation: %w", err)
	}
	r.Get("/", valsHandler.MetricValuesHandler())

	return nil
}

func addPingHandlers(r chi.Router, s handlers.Service, log *zap.Logger) error {
	pingHandler, err := handlers.NewPingHandler(s, log)
	if err != nil {
		return fmt.Errorf("ping handler creation: %w", err)
	}
	r.Get("/ping", pingHandler.Handler())

	return nil
}
