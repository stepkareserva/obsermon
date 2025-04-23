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

	updateHandler, err := updateHandler(ctx, s, log)
	if err != nil {
		return nil, err
	}
	valueHandler, err := valueHandler(ctx, s, log)
	if err != nil {
		return nil, err
	}

	valuesHandler, err := valuesHandler(ctx, s, log)
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()
	if log != nil {
		r.Use(middleware.Logger(log))
	}
	r.Use(middleware.Compression(log))
	r.Use(middleware.Buffering(log))

	r.Mount("/update", updateHandler)
	r.Mount("/value", valueHandler)
	r.Mount("/", valuesHandler)

	return r, nil
}

func updateHandler(ctx context.Context, s Service, log *zap.Logger) (http.Handler, error) {
	if s == nil {
		return nil, fmt.Errorf("metrics service is nil")
	}
	handler, err := NewUpdateHandler(s, log)
	if err != nil {
		return nil, fmt.Errorf("value handler creation: %w", err)
	}
	r := chi.NewRouter()

	r.Post(fmt.Sprintf("/%s/{%s}/{%s}", MetricGauge, ChiName, ChiValue),
		handler.UpdateGaugeURLHandler(ctx))
	r.Post(fmt.Sprintf("/%s/{%s}/{%s}", MetricCounter, ChiName, ChiValue),
		handler.UpdateCounterURLHandler(ctx))
	r.Post(fmt.Sprintf("/{%s}/{%s}/{%s}", ChiMetric, ChiName, ChiValue),
		handler.UpdateUnknownMetricURLHandler(ctx))
	r.Post("/",
		handler.UpdateMetricJSONHandler(ctx))

	return r, nil
}

func valueHandler(ctx context.Context, s Service, log *zap.Logger) (http.Handler, error) {
	if s == nil {
		return nil, fmt.Errorf("metrics server is nil")
	}
	handler, err := NewValueHandler(s, log)
	if err != nil {
		return nil, fmt.Errorf("value handler creation: %w", err)
	}

	r := chi.NewRouter()

	r.Get(fmt.Sprintf("/%s/{%s}", MetricGauge, ChiName),
		handler.GaugeValueURLHandler(ctx))
	r.Get(fmt.Sprintf("/%s/{%s}", MetricCounter, ChiName),
		handler.CounterValueURLHandler(ctx))
	r.Get(fmt.Sprintf("/{%s}/{%s}", ChiMetric, ChiName),
		handler.UnknownMetricValueURLHandler(ctx))
	r.Post("/",
		handler.ValueMetricJSONHandler(ctx))

	return r, nil
}

func valuesHandler(ctx context.Context, s Service, log *zap.Logger) (http.Handler, error) {
	if s == nil {
		return nil, fmt.Errorf("metrics service is nil")
	}
	handler, err := NewValuesHandler(s, log)
	if err != nil {
		return nil, fmt.Errorf("values handler creation: %w", err)
	}

	r := chi.NewRouter()
	r.Get("/",
		handler.MetricValuesHandler(ctx))

	return r, nil
}
