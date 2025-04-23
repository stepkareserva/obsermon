package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/stepkareserva/obsermon/internal/server/middleware"
	"go.uber.org/zap"
)

type MetricsHandler struct {
	service Service
	log     *zap.Logger
}

func New(s Service, log *zap.Logger) (*MetricsHandler, error) {
	if s == nil {
		return nil, fmt.Errorf("service not exist")
	}
	return &MetricsHandler{
		service: s,
		log:     log,
	}, nil
}

func (h *MetricsHandler) Handler(ctx context.Context) (http.Handler, error) {
	if h == nil {
		return nil, fmt.Errorf("metrics handler not exists")
	}
	if h.service == nil {
		return nil, fmt.Errorf("metrics service is nil")
	}

	updateHandler, err := h.updateHandler(ctx)
	if err != nil {
		return nil, err
	}
	valueHandler, err := h.valueHandler(ctx)
	if err != nil {
		return nil, err
	}

	valuesHandler, err := h.valuesHandler(ctx)
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()
	if h.log != nil {
		r.Use(middleware.Logger(h.log))
	}
	r.Use(middleware.Compression(h.log))
	r.Use(middleware.Buffering(h.log))

	r.Mount("/update", updateHandler)
	r.Mount("/value", valueHandler)
	r.Mount("/", valuesHandler)

	return r, nil
}

func (h *MetricsHandler) updateHandler(ctx context.Context) (http.Handler, error) {
	handler, err := NewUpdateHandler(h.service, h.log)
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

func (h *MetricsHandler) valueHandler(ctx context.Context) (http.Handler, error) {
	handler, err := NewValueHandler(h.service, h.log)
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

func (h *MetricsHandler) valuesHandler(ctx context.Context) (http.Handler, error) {
	handler, err := NewValuesHandler(h.service, h.log)
	if err != nil {
		return nil, fmt.Errorf("values handler creation: %w", err)
	}

	r := chi.NewRouter()
	r.Get("/",
		handler.MetricValuesHandler(ctx))

	return r, nil
}
