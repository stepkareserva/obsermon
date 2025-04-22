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

	r := chi.NewRouter()

	r.Post(fmt.Sprintf("/%s/{%s}/{%s}", MetricGauge, ChiName, ChiValue),
		updateGaugeURLHandler(ctx, s, log))
	r.Post(fmt.Sprintf("/%s/{%s}/{%s}", MetricCounter, ChiName, ChiValue),
		updateCounterURLHandler(ctx, s, log))
	r.Post(fmt.Sprintf("/{%s}/{%s}/{%s}", ChiMetric, ChiName, ChiValue),
		updateUnknownMetricURLHandler(ctx, log))
	r.Post("/", updateMetricJSONHandler(ctx, s, log))

	return r, nil
}

func valueHandler(ctx context.Context, s Service, log *zap.Logger) (http.Handler, error) {
	if s == nil {
		return nil, fmt.Errorf("metrics server is nil")
	}

	r := chi.NewRouter()

	r.Get(fmt.Sprintf("/%s/{%s}", MetricGauge, ChiName),
		gaugeValueURLHandler(ctx, s, log))
	r.Get(fmt.Sprintf("/%s/{%s}", MetricCounter, ChiName),
		counterValueURLHandler(ctx, s, log))
	r.Get(fmt.Sprintf("/{%s}/{%s}", ChiMetric, ChiName),
		unknownMetricValueURLHandler(ctx, log))
	r.Post("/", valueMetricJSONHandler(ctx, s, log))

	return r, nil
}

func valuesHandler(ctx context.Context, s Service, log *zap.Logger) (http.Handler, error) {
	if s == nil {
		return nil, fmt.Errorf("metrics service is nil")
	}

	r := chi.NewRouter()
	r.Get("/", metricValuesHandler(ctx, s, log))

	return r, nil
}
