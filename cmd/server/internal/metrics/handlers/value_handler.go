package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/stepkareserva/obsermon/cmd/server/internal/metrics/service"
)

func ValueHandler(s *service.Service) (http.Handler, error) {
	if s == nil {
		return nil, fmt.Errorf("metrics server is nil")
	}

	r := chi.NewRouter()

	r.Get(fmt.Sprintf("/%s/{%s}", MetricGauge, ChiName),
		gaugeValueHandler(s))
	r.Get(fmt.Sprintf("/%s/{%s}", MetricCounter, ChiName),
		counterValueHandler(s))
	r.Get(fmt.Sprintf("/{%s}/{%s}", ChiMetric, ChiName),
		unknownMetricValueHandler())

	return r, nil
}

func gaugeValueHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, ChiName)
		gauge, exists, err := s.GetGauge(name)
		if err != nil {
			WriteError(w, ErrInternalServerError)
			return
		}

		if !exists {
			WriteError(w, ErrMetricNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		if _, err := w.Write([]byte(gauge.Value.String())); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func counterValueHandler(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, ChiName)
		counter, exists, err := s.GetCounter(name)
		if err != nil {
			WriteError(w, ErrInternalServerError)
			return
		}

		if !exists {
			WriteError(w, ErrMetricNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if _, err := w.Write([]byte(counter.Value.String())); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func unknownMetricValueHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteError(w, ErrInvalidMetricType, chi.URLParam(r, ChiMetric))
	}
}
