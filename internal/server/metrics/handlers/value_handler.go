package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/go-playground/validator"
	"github.com/stepkareserva/obsermon/internal/models"

	hc "github.com/stepkareserva/obsermon/internal/server/httpconst"
)

func gaugeValueURLHandler(ctx context.Context, s Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, ChiName)
		gauge, exists, err := s.GetGauge(name)
		if err != nil {
			WriteError(w, ErrInternalServerError, log)
			return
		}

		if !exists {
			WriteError(w, ErrMetricNotFound, log)
			return
		}

		w.Header().Set(hc.ContentType, hc.ContentTypeText)
		if _, err := w.Write([]byte(gauge.Value.String())); err != nil {
			WriteError(w, ErrInternalServerError, log)
			return
		}
	}
}

func counterValueURLHandler(ctx context.Context, s Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, ChiName)
		counter, exists, err := s.GetCounter(name)
		if err != nil {
			WriteError(w, ErrInternalServerError, log)
			return
		}

		if !exists {
			WriteError(w, ErrMetricNotFound, log)
			return
		}

		w.Header().Set(hc.ContentType, hc.ContentTypeText)
		if _, err := w.Write([]byte(counter.Value.String())); err != nil {
			WriteError(w, ErrInternalServerError, log)
			return
		}
	}
}

func unknownMetricValueURLHandler(ctx context.Context, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteError(w, ErrInvalidMetricType, log, chi.URLParam(r, ChiMetric))
	}
}

func valueMetricJSONHandler(ctx context.Context, s Service, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(hc.ContentType) != hc.ContentTypeJSON {
			WriteError(w, ErrUnsupportedContentType, log)
			return
		}
		var request models.MetricValueRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			WriteError(w, ErrInvalidRequestJSON, log)
			return
		}
		if err := validator.New().Struct(request); err != nil {
			WriteError(w, ErrInvalidRequestJSON, log)
			return
		}
		m, exists, err := s.GetMetric(request.MType, request.ID)
		if err != nil {
			WriteError(w, ErrInternalServerError, log)
			return
		}
		if !exists {
			WriteError(w, ErrMetricNotFound, log)
			return
		}
		w.Header().Set(hc.ContentType, hc.ContentTypeJSON)
		if err = json.NewEncoder(w).Encode(m); err != nil {
			WriteError(w, ErrInternalServerError, log)
			return
		}
	}
}
