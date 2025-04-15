package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/go-playground/validator"
	"github.com/stepkareserva/obsermon/internal/models"

	hc "github.com/stepkareserva/obsermon/internal/server/httpconst"
)

func ValueHandler(s Service) (http.Handler, error) {
	if s == nil {
		return nil, fmt.Errorf("metrics server is nil")
	}

	r := chi.NewRouter()

	r.Get(fmt.Sprintf("/%s/{%s}", MetricGauge, ChiName),
		gaugeValueURLHandler(s))
	r.Get(fmt.Sprintf("/%s/{%s}", MetricCounter, ChiName),
		counterValueURLHandler(s))
	r.Get(fmt.Sprintf("/{%s}/{%s}", ChiMetric, ChiName),
		unknownMetricValueURLHandler())
	r.Post("/", valueMetricJSONHandler(s))

	return r, nil
}

func gaugeValueURLHandler(s Service) http.HandlerFunc {
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

		w.Header().Set(hc.ContentType, hc.ContentTypeText)
		if _, err := w.Write([]byte(gauge.Value.String())); err != nil {
			WriteError(w, ErrInternalServerError)
			return
		}
	}
}

func counterValueURLHandler(s Service) http.HandlerFunc {
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

		w.Header().Set(hc.ContentType, hc.ContentTypeText)
		if _, err := w.Write([]byte(counter.Value.String())); err != nil {
			WriteError(w, ErrInternalServerError)
			return
		}
	}
}

func unknownMetricValueURLHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteError(w, ErrInvalidMetricType, chi.URLParam(r, ChiMetric))
	}
}

func valueMetricJSONHandler(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(hc.ContentType) != hc.ContentTypeJSON {
			WriteError(w, ErrUnsupportedContentType)
			return
		}
		var request models.MetricValueRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			WriteError(w, ErrInvalidRequestJSON)
			return
		}
		if err := validator.New().Struct(request); err != nil {
			WriteError(w, ErrInvalidRequestJSON)
			return
		}
		m, exists, err := s.GetMetric(request.MType, request.ID)
		if err != nil {
			WriteError(w, ErrInternalServerError)
			return
		}
		if !exists {
			WriteError(w, ErrMetricNotFound)
			return
		}
		w.Header().Set(hc.ContentType, hc.ContentTypeJSON)
		if err = json.NewEncoder(w).Encode(m); err != nil {
			WriteError(w, ErrInternalServerError)
			return
		}
	}
}
