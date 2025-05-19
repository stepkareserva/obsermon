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
	"github.com/stepkareserva/obsermon/internal/server/http/constants"
	"github.com/stepkareserva/obsermon/internal/server/http/errors"
)

type ValueHandler struct {
	service Service
	errors.ErrorsWriter
}

func NewValueHandler(s Service, log *zap.Logger) (*ValueHandler, error) {
	if s == nil {
		return nil, fmt.Errorf("service not exists")
	}
	return &ValueHandler{
		service:      s,
		ErrorsWriter: errors.NewErrorsWriter(log),
	}, nil
}

func (h *ValueHandler) GaugeValueURLHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, constants.ChiName)
		gauge, exists, err := h.service.FindGauge(r.Context(), name)
		if err != nil {
			h.WriteError(w, errors.ErrInternalServerError, err.Error())
			return
		}

		if !exists {
			h.WriteError(w, errors.ErrMetricNotFound)
			return
		}

		w.Header().Set(constants.ContentType, constants.ContentTypeText)
		if _, err := w.Write([]byte(gauge.Value.String())); err != nil {
			h.WriteError(w, errors.ErrInternalServerError, err.Error())
			return
		}
	}
}

func (h *ValueHandler) CounterValueURLHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, constants.ChiName)
		counter, exists, err := h.service.FindCounter(r.Context(), name)
		if err != nil {
			h.WriteError(w, errors.ErrInternalServerError, err.Error())
			return
		}

		if !exists {
			h.WriteError(w, errors.ErrMetricNotFound)
			return
		}

		w.Header().Set(constants.ContentType, constants.ContentTypeText)
		if _, err := w.Write([]byte(counter.Value.String())); err != nil {
			h.WriteError(w, errors.ErrInternalServerError, err.Error())
			return
		}
	}
}

func (h *ValueHandler) UnknownMetricValueURLHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.WriteError(w, errors.ErrInvalidMetricType, chi.URLParam(r, constants.ChiMetric))
	}
}

func (h *ValueHandler) ValueMetricJSONHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(constants.ContentType) != constants.ContentTypeJSON {
			h.WriteError(w, errors.ErrUnsupportedContentType)
			return
		}
		var request models.MetricValueRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			h.WriteError(w, errors.ErrInvalidRequestJSON, err.Error())
			return
		}
		if err := validator.New().Struct(request); err != nil {
			h.WriteError(w, errors.ErrInvalidRequestJSON, err.Error())
			return
		}
		m, exists, err := h.service.FindMetric(r.Context(), request.MType, request.ID)
		if err != nil {
			h.WriteError(w, errors.ErrInternalServerError, err.Error())
			return
		}
		if !exists {
			h.WriteError(w, errors.ErrMetricNotFound)
			return
		}
		w.Header().Set(constants.ContentType, constants.ContentTypeJSON)
		if err = json.NewEncoder(w).Encode(m); err != nil {
			h.WriteError(w, errors.ErrInternalServerError, err.Error())
			return
		}
	}
}
