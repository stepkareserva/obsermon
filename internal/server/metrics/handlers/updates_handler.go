package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/go-playground/validator"
	"github.com/stepkareserva/obsermon/internal/models"

	hu "github.com/stepkareserva/obsermon/internal/server/httputils"
)

type UpdatesHandler struct {
	service Service
	hu.ErrorsWriter
}

func NewUpdatesHandler(s Service, log *zap.Logger) (*UpdatesHandler, error) {
	if s == nil {
		return nil, fmt.Errorf("service not exists")
	}
	return &UpdatesHandler{
		service:      s,
		ErrorsWriter: hu.NewErrorsWriter(log),
	}, nil
}

func (h *UpdateHandler) UpdateMetricsJSONHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(hu.ContentType) != hu.ContentTypeJSON {
			h.WriteError(w, hu.ErrUnsupportedContentType)
			return
		}
		var request models.UpdateMetricsRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			h.WriteError(w, hu.ErrInvalidRequestJSON, err.Error())
			return
		}

		v := validator.New()
		for _, requestItem := range request {
			if err := v.Struct(requestItem); err != nil {
				h.WriteError(w, hu.ErrInvalidRequestJSON, err.Error())
				return
			}
		}
		updated, err := h.service.UpdateMetrics(request)
		if err != nil {
			h.WriteError(w, hu.ErrInternalServerError, err.Error())
			return
		}
		w.Header().Set(hu.ContentType, hu.ContentTypeJSON)
		if err = json.NewEncoder(w).Encode(updated); err != nil {
			h.WriteError(w, hu.ErrInternalServerError, err.Error())
			return
		}
	}
}
