package handlers

import (
	"net/http"

	"github.com/stepkareserva/obsermon/internal/server/http/constants"
	"github.com/stepkareserva/obsermon/internal/server/http/errors"
	"go.uber.org/zap"
)

type PingHandler struct {
	service Service
	errors.ErrorsWriter
}

func NewPingHandler(service Service, log *zap.Logger) (*PingHandler, error) {
	return &PingHandler{
		service:      service,
		ErrorsWriter: errors.NewErrorsWriter(log),
	}, nil
}

func (h *PingHandler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h.service.Ping(r.Context()); err != nil {
			h.ErrorsWriter.WriteError(w, errors.ErrDatabaseUnavailable, err.Error())
			return
		}

		w.Header().Set(constants.ContentType, constants.ContentTypeText)
		w.WriteHeader(http.StatusOK)
	}
}
