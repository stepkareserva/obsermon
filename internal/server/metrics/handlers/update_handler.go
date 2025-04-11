package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"

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
		updateGaugeURLHandler(s))
	r.Post(fmt.Sprintf("/%s/{%s}/{%s}", MetricCounter, ChiName, ChiValue),
		updateCounterURLHandler(s))
	r.Post(fmt.Sprintf("/{%s}/{%s}/{%s}", ChiMetric, ChiName, ChiValue),
		updateUnknownMetricURLHandler())
	r.Post("/", updateMetricJSONHandler(s))

	return r, nil
}

func updateGaugeURLHandler(s Service) http.HandlerFunc {
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

func updateCounterURLHandler(s Service) http.HandlerFunc {
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

func updateUnknownMetricURLHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteError(w, ErrInvalidMetricType, chi.URLParam(r, ChiMetric))
	}
}

func updateMetricJSONHandler(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(contentType) != contentTypeJSON {
			WriteError(w, ErrUnsupportedContentType)
			return
		}
		var request models.UpdateMetricRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			WriteError(w, ErrInvalidRequestJSON)
			return
		}
		if err := validator.New().Struct(request); err != nil {
			WriteError(w, ErrInvalidRequestJSON)
			return
		}
		if err := s.UpdateMetric(request); err != nil {
			WriteError(w, ErrInternalServerError)
			return
		}
		m, exists, err := s.GetMetric(request.MType, request.ID)
		if err != nil || !exists {
			WriteError(w, ErrInternalServerError)
			return
		}
		w.Header().Set(contentType, contentTypeJSON)
		if err = json.NewEncoder(w).Encode(m); err != nil {
			WriteError(w, ErrInternalServerError)
			return
		}
	}
}
