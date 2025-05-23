package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"go.uber.org/zap"

	"github.com/stepkareserva/obsermon/internal/models"

	"github.com/stepkareserva/obsermon/internal/server/http/constants"
	"github.com/stepkareserva/obsermon/internal/server/http/errors"
)

type UpdateHandler struct {
	service Service
	errors.ErrorsWriter
}

func NewUpdateHandler(s Service, log *zap.Logger) (*UpdateHandler, error) {
	if s == nil {
		return nil, fmt.Errorf("service not exists")
	}
	return &UpdateHandler{
		service:      s,
		ErrorsWriter: errors.NewErrorsWriter(log),
	}, nil
}

func (h *UpdateHandler) UpdateGaugeURLHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, constants.ChiName)
		var value models.GaugeValue
		if err := value.FromString(chi.URLParam(r, constants.ChiValue)); err != nil {
			h.WriteError(w, errors.ErrInvalidMetricValue, err.Error())
			return
		}
		gauge := models.Gauge{Name: name, Value: value}
		if _, err := h.service.UpdateGauge(r.Context(), gauge); err != nil {
			h.WriteError(w, errors.ErrInternalServerError, err.Error())
			return
		}

		w.Header().Set(constants.ContentType, constants.ContentTypeText)
		w.WriteHeader(http.StatusOK)
	}
}

func (h *UpdateHandler) UpdateCounterURLHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, constants.ChiName)
		var value models.CounterValue
		if err := value.FromString(chi.URLParam(r, constants.ChiValue)); err != nil {
			h.WriteError(w, errors.ErrInvalidMetricValue, err.Error())
			return
		}

		counter := models.Counter{Name: name, Value: value}
		if _, err := h.service.UpdateCounter(r.Context(), counter); err != nil {
			h.WriteError(w, errors.ErrInternalServerError, err.Error())
			return
		}

		w.Header().Set(constants.ContentType, constants.ContentTypeHTML)
		w.WriteHeader(http.StatusOK)
	}
}

func (h *UpdateHandler) UpdateUnknownMetricURLHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.WriteError(w, errors.ErrInvalidMetricType, chi.URLParam(r, constants.ChiMetric))
	}
}

func (h *UpdateHandler) UpdateMetricJSONHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(constants.ContentType) != constants.ContentTypeJSON {
			h.WriteError(w, errors.ErrUnsupportedContentType)
			return
		}
		var request models.UpdateMetricRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			h.WriteError(w, errors.ErrInvalidRequestJSON, err.Error())
			return
		}
		if err := validator.New().Struct(request); err != nil {
			h.WriteError(w, errors.ErrInvalidRequestJSON, err.Error())
			return
		}
		updated, err := h.service.UpdateMetric(r.Context(), request)
		if err != nil {
			h.WriteError(w, errors.ErrInternalServerError, err.Error())
			return
		}
		// update and return updated metrics in the same request
		// may be bottleneck in scenarious where updates are
		// frequent and value requests are rate.
		// in such scenarious we can only collect metrics on
		// update and aggregate on value requests.
		// but now updated metric value should be returned as
		// part of update request's response, so it keep in mind.
		w.Header().Set(constants.ContentType, constants.ContentTypeJSON)
		if err = json.NewEncoder(w).Encode(*updated); err != nil {
			h.WriteError(w, errors.ErrInternalServerError, err.Error())
			return
		}
	}
}
