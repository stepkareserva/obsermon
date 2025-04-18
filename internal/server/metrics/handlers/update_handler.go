package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
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

func UpdateHandler(s Service) (http.Handler, error) {
	if s == nil {
		return nil, fmt.Errorf("metrics service is nil")
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

func updateGaugeHandler(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, ChiName)
		var value models.GaugeValue
		if err := value.FromString(chi.URLParam(r, ChiValue)); err != nil {
			WriteError(w, ErrInvalidMetricValue)
			return
		}
		gauge := models.Gauge{Name: name, Value: value}
		if err := s.UpdateGauge(gauge); err != nil {
			WriteError(w, ErrInternalServerError)
			return
		}

		w.Header().Set(contentType, contentTypeText)
		w.WriteHeader(http.StatusOK)
	}
}

func updateCounterHandler(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, ChiName)
		var value models.CounterValue
		if err := value.FromString(chi.URLParam(r, ChiValue)); err != nil {
			WriteError(w, ErrInvalidMetricValue)
			return
		}

		counter := models.Counter{Name: name, Value: value}
		if err := s.UpdateCounter(counter); err != nil {
			WriteError(w, ErrInternalServerError)
			return
		}

		w.Header().Set(contentType, contentTypeHTML)
		w.WriteHeader(http.StatusOK)
	}
}

func updateUnknownMetricHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteError(w, ErrInvalidMetricType, chi.URLParam(r, ChiMetric))
	}
}
