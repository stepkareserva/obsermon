package handlers

import (
	"context"
	"fmt"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handlers struct {
	updHandler  *UpdateHandler
	valHandler  *ValueHandler
	valsHandler *ValuesHandler
}

func New(s Service, log *zap.Logger) (*Handlers, error) {
	if s == nil {
		return nil, fmt.Errorf("service not exist")
	}

	updHandler, err := NewUpdateHandler(s, log)
	if err != nil {
		return nil, fmt.Errorf("update handler creation: %w", err)
	}
	valHandler, err := NewValueHandler(s, log)
	if err != nil {
		return nil, fmt.Errorf("value handler creation: %w", err)
	}
	valsHandler, err := NewValuesHandler(s, log)
	if err != nil {
		return nil, fmt.Errorf("values handler creation: %w", err)
	}

	return &Handlers{
		updHandler:  updHandler,
		valHandler:  valHandler,
		valsHandler: valsHandler,
	}, nil

}

func (h *Handlers) RegisterRoutes(ctx context.Context, r chi.Router) error {
	if h == nil {
		return fmt.Errorf("handlers not exist")
	}

	r.Route("/update", func(r chi.Router) {
		h.registerUpdateRoutes(ctx, r)
	})
	r.Route("/value", func(r chi.Router) {
		h.registerValueRoutes(ctx, r)
	})
	r.Route("/", func(r chi.Router) {
		h.registerValuesRoutes(ctx, r)
	})

	return nil
}

func (h *Handlers) registerUpdateRoutes(ctx context.Context, r chi.Router) {
	r.Post(fmt.Sprintf("/%s/{%s}/{%s}", MetricGauge, ChiName, ChiValue),
		h.updHandler.UpdateGaugeURLHandler(ctx))
	r.Post(fmt.Sprintf("/%s/{%s}/{%s}", MetricCounter, ChiName, ChiValue),
		h.updHandler.UpdateCounterURLHandler(ctx))
	r.Post(fmt.Sprintf("/{%s}/{%s}/{%s}", ChiMetric, ChiName, ChiValue),
		h.updHandler.UpdateUnknownMetricURLHandler(ctx))
	r.Post("/",
		h.updHandler.UpdateMetricJSONHandler(ctx))
}

func (h *Handlers) registerValueRoutes(ctx context.Context, r chi.Router) {
	r.Get(fmt.Sprintf("/%s/{%s}", MetricGauge, ChiName),
		h.valHandler.GaugeValueURLHandler(ctx))
	r.Get(fmt.Sprintf("/%s/{%s}", MetricCounter, ChiName),
		h.valHandler.CounterValueURLHandler(ctx))
	r.Get(fmt.Sprintf("/{%s}/{%s}", ChiMetric, ChiName),
		h.valHandler.UnknownMetricValueURLHandler(ctx))
	r.Post("/",
		h.valHandler.ValueMetricJSONHandler(ctx))
}

func (h *Handlers) registerValuesRoutes(ctx context.Context, r chi.Router) {
	r.Get("/",
		h.valsHandler.MetricValuesHandler(ctx))
}
