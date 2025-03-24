package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/stepkareserva/obsermon/cmd/server/internal/metrics/server"
	"github.com/stepkareserva/obsermon/internal/models"
)

const (
	// metrics for url
	MetricGauge   = "gauge"
	MetricCounter = "counter"

	// names of chi routing url params to be extracted
	ChiMetric = "metric"
	ChiName   = "name"
	ChiValue  = "value"
)

func UpdateHandler(s *server.Server) (http.Handler, error) {
	if s == nil {
		return nil, fmt.Errorf("metrics server is nil")
	}

	r := chi.NewRouter()

	r.Post(fmt.Sprintf("/%s/{%s}/{%s}", MetricGauge, ChiName, ChiValue),
		updateGaugeHandler(s))
	r.Post(fmt.Sprintf("/%s/{%s}/{%s}", MetricCounter, ChiName, ChiValue),
		updateCounterHandler(s))
	r.Post(fmt.Sprintf("/{%s}/{%s}/{%s}", ChiMetric, ChiName, ChiValue),
		updateUnknownMetricHandler())

	return r, nil
}

func updateGaugeHandler(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, ChiName)
		var value models.Gauge
		if err := value.FromString(chi.URLParam(r, ChiValue)); err != nil {
			WriteError(w, ErrInvalidMetricValue)
			return
		}

		if err := s.UpdateGauge(name, value); err != nil {
			WriteError(w, ErrInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
	}
}

func updateCounterHandler(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, ChiName)
		var value models.Counter
		if err := value.FromString(chi.URLParam(r, ChiValue)); err != nil {
			WriteError(w, ErrInvalidMetricValue)
			return
		}

		if err := s.UpdateCounter(name, value); err != nil {
			WriteError(w, ErrInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
	}
}

func updateUnknownMetricHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteError(w, ErrInvalidMetricType, chi.URLParam(r, ChiMetric))
	}
}
