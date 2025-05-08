package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/go-playground/validator"
	"github.com/stepkareserva/obsermon/internal/models"

	hu "github.com/stepkareserva/obsermon/internal/server/httputils"
)

type ValueHandler struct {
	service Service
	hu.ErrorsWriter
}

func NewValueHandler(s Service, log *zap.Logger) (*ValueHandler, error) {
	if s == nil {
		return nil, fmt.Errorf("service not exists")
	}
	return &ValueHandler{
		service:      s,
		ErrorsWriter: hu.NewErrorsWriter(log),
	}, nil
}

func (h *ValueHandler) GaugeValueURLHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, ChiName)
		gauge, exists, err := h.service.FindGauge(r.Context(), name)
		if err != nil {
			h.WriteError(w, hu.ErrInternalServerError, err.Error())
			return
		}

		if !exists {
			h.WriteError(w, ErrMetricNotFound)
			return
		}

		w.Header().Set(hu.ContentType, hu.ContentTypeText)
		if _, err := w.Write([]byte(gauge.Value.String())); err != nil {
			h.WriteError(w, hu.ErrInternalServerError, err.Error())
			return
		}
	}
}

func (h *ValueHandler) CounterValueURLHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, ChiName)
		counter, exists, err := h.service.FindCounter(r.Context(), name)
		if err != nil {
			h.WriteError(w, hu.ErrInternalServerError, err.Error())
			return
		}

		if !exists {
			h.WriteError(w, ErrMetricNotFound)
			return
		}

		w.Header().Set(hu.ContentType, hu.ContentTypeText)
		if _, err := w.Write([]byte(counter.Value.String())); err != nil {
			h.WriteError(w, hu.ErrInternalServerError, err.Error())
			return
		}
	}
}

func (h *ValueHandler) UnknownMetricValueURLHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.WriteError(w, ErrInvalidMetricType, chi.URLParam(r, ChiMetric))
	}
}

func (h *ValueHandler) ValueMetricJSONHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(hu.ContentType) != hu.ContentTypeJSON {
			h.WriteError(w, hu.ErrUnsupportedContentType)
			return
		}
		var request models.MetricValueRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			h.WriteError(w, hu.ErrInvalidRequestJSON, err.Error())
			return
		}
		if err := validator.New().Struct(request); err != nil {
			h.WriteError(w, hu.ErrInvalidRequestJSON, err.Error())
			return
		}
		m, exists, err := h.service.FindMetric(r.Context(), request.MType, request.ID)
		if err != nil {
			h.WriteError(w, hu.ErrInternalServerError, err.Error())
			return
		}
		if !exists {
			h.WriteError(w, ErrMetricNotFound)
			return
		}
		w.Header().Set(hu.ContentType, hu.ContentTypeJSON)
		if err = json.NewEncoder(w).Encode(m); err != nil {
			h.WriteError(w, hu.ErrInternalServerError, err.Error())
			return
		}
	}
}
