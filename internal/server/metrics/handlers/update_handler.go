package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"go.uber.org/zap"

	"github.com/stepkareserva/obsermon/internal/models"

	hc "github.com/stepkareserva/obsermon/internal/server/httpconst"
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

func updateGaugeURLHandler(s Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, ChiName)
		var value models.GaugeValue
		if err := value.FromString(chi.URLParam(r, ChiValue)); err != nil {
			WriteError(w, ErrInvalidMetricValue, log)
			return
		}
		gauge := models.Gauge{Name: name, Value: value}
		if err := s.UpdateGauge(gauge); err != nil {
			WriteError(w, ErrInternalServerError, log)
			return
		}

		w.Header().Set(hc.ContentType, hc.ContentTypeText)
		w.WriteHeader(http.StatusOK)
	}
}

func updateCounterURLHandler(s Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, ChiName)
		var value models.CounterValue
		if err := value.FromString(chi.URLParam(r, ChiValue)); err != nil {
			WriteError(w, ErrInvalidMetricValue, log)
			return
		}

		counter := models.Counter{Name: name, Value: value}
		if err := s.UpdateCounter(counter); err != nil {
			WriteError(w, ErrInternalServerError, log)
			return
		}

		w.Header().Set(hc.ContentType, hc.ContentTypeHTML)
		w.WriteHeader(http.StatusOK)
	}
}

func updateUnknownMetricURLHandler(log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteError(w, ErrInvalidMetricType, log, chi.URLParam(r, ChiMetric))
	}
}

func updateMetricJSONHandler(s Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(hc.ContentType) != hc.ContentTypeJSON {
			WriteError(w, ErrUnsupportedContentType, log)
			return
		}
		var request models.UpdateMetricRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			WriteError(w, ErrInvalidRequestJSON, log)
			return
		}
		if err := validator.New().Struct(request); err != nil {
			WriteError(w, ErrInvalidRequestJSON, log)
			return
		}
		if err := s.UpdateMetric(request); err != nil {
			WriteError(w, ErrInternalServerError, log)
			return
		}
		// update and return updated metrics in the same request
		// may be bottleneck in scenarious where updates are
		// frequent and value requests are rate.
		// in such scenarious we can only collect metrics on
		// update and aggregate on value requests.
		// but now updated metric value should be returned as
		// part of update request's response, so it keep in mind.
		// for better performance methods UpdateMetric and
		// UpdateAndExtractMetric required, i think, but...
		m, exists, err := s.GetMetric(request.MType, request.ID)
		if err != nil || !exists {
			WriteError(w, ErrInternalServerError, log)
			return
		}
		w.Header().Set(hc.ContentType, hc.ContentTypeJSON)
		if err = json.NewEncoder(w).Encode(m); err != nil {
			WriteError(w, ErrInternalServerError, log)
			return
		}
	}
}
