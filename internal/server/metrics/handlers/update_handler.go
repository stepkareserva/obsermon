package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"go.uber.org/zap"

	"github.com/stepkareserva/obsermon/internal/models"

	hu "github.com/stepkareserva/obsermon/internal/server/httputils"
)

type UpdateHandler struct {
	service Service
	hu.ErrorsWriter
}

func NewUpdateHandler(s Service, log *zap.Logger) (*UpdateHandler, error) {
	if s == nil {
		return nil, fmt.Errorf("service not exists")
	}
	return &UpdateHandler{
		service:      s,
		ErrorsWriter: hu.NewErrorsWriter(log),
	}, nil
}

func (h *UpdateHandler) UpdateGaugeURLHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, ChiName)
		var value models.GaugeValue
		if err := value.FromString(chi.URLParam(r, ChiValue)); err != nil {
			h.WriteError(w, ErrInvalidMetricValue)
			return
		}
		gauge := models.Gauge{Name: name, Value: value}
		if _, err := h.service.UpdateGauge(gauge); err != nil {
			h.WriteError(w, hu.ErrInternalServerError)
			return
		}

		w.Header().Set(hu.ContentType, hu.ContentTypeText)
		w.WriteHeader(http.StatusOK)
	}
}

func (h *UpdateHandler) UpdateCounterURLHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, ChiName)
		var value models.CounterValue
		if err := value.FromString(chi.URLParam(r, ChiValue)); err != nil {
			h.WriteError(w, ErrInvalidMetricValue)
			return
		}

		counter := models.Counter{Name: name, Value: value}
		if _, err := h.service.UpdateCounter(counter); err != nil {
			h.WriteError(w, hu.ErrInternalServerError)
			return
		}

		w.Header().Set(hu.ContentType, hu.ContentTypeHTML)
		w.WriteHeader(http.StatusOK)
	}
}

func (h *UpdateHandler) UpdateUnknownMetricURLHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.WriteError(w, ErrInvalidMetricType, chi.URLParam(r, ChiMetric))
	}
}

func (h *UpdateHandler) UpdateMetricJSONHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(hu.ContentType) != hu.ContentTypeJSON {
			h.WriteError(w, hu.ErrUnsupportedContentType)
			return
		}
		var request models.UpdateMetricRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			h.WriteError(w, hu.ErrInvalidRequestJSON)
			return
		}
		if err := validator.New().Struct(request); err != nil {
			h.WriteError(w, hu.ErrInvalidRequestJSON)
			return
		}
		updated, err := h.service.UpdateMetric(request)
		if err != nil {
			h.WriteError(w, hu.ErrInternalServerError)
			return
		}
		// update and return updated metrics in the same request
		// may be bottleneck in scenarious where updates are
		// frequent and value requests are rate.
		// in such scenarious we can only collect metrics on
		// update and aggregate on value requests.
		// but now updated metric value should be returned as
		// part of update request's response, so it keep in mind.
		w.Header().Set(hu.ContentType, hu.ContentTypeJSON)
		if err = json.NewEncoder(w).Encode(*updated); err != nil {
			h.WriteError(w, hu.ErrInternalServerError)
			return
		}
	}
}
