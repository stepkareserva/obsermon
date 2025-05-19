package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/go-playground/validator"
	"github.com/stepkareserva/obsermon/internal/models"

	"github.com/stepkareserva/obsermon/internal/server/http/constants"
	"github.com/stepkareserva/obsermon/internal/server/http/errors"
)

type UpdatesHandler struct {
	service Service
	errors.ErrorsWriter
}

func NewUpdatesHandler(s Service, log *zap.Logger) (*UpdatesHandler, error) {
	if s == nil {
		return nil, fmt.Errorf("service not exists")
	}
	return &UpdatesHandler{
		service:      s,
		ErrorsWriter: errors.NewErrorsWriter(log),
	}, nil
}

func (h *UpdateHandler) UpdateMetricsJSONHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(constants.ContentType) != constants.ContentTypeJSON {
			h.WriteError(w, errors.ErrUnsupportedContentType)
			return
		}
		var request models.UpdateMetricsRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			h.WriteError(w, errors.ErrInvalidRequestJSON, err.Error())
			return
		}

		v := validator.New()
		for _, requestItem := range request {
			if err := v.Struct(requestItem); err != nil {
				h.WriteError(w, errors.ErrInvalidRequestJSON, err.Error())
				return
			}
		}
		updated, err := h.service.UpdateMetrics(r.Context(), request)
		if err != nil {
			h.WriteError(w, errors.ErrInternalServerError, err.Error())
			return
		}
		w.Header().Set(constants.ContentType, constants.ContentTypeJSON)
		if err = json.NewEncoder(w).Encode(updated); err != nil {
			h.WriteError(w, errors.ErrInternalServerError, err.Error())
			return
		}
	}
}
