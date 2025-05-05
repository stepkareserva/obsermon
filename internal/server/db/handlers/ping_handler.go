package handlers

import (
	"context"
	"net/http"

	hu "github.com/stepkareserva/obsermon/internal/server/httputils"
	"go.uber.org/zap"
)

type PingHandler struct {
	db Database
	hu.ErrorsWriter
}

func NewPingHandler(db Database, log *zap.Logger) (*PingHandler, error) {
	return &PingHandler{
		db:           db,
		ErrorsWriter: hu.NewErrorsWriter(log),
	}, nil
}

func (h *PingHandler) Handler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h.db.Ping(); err != nil {
			h.ErrorsWriter.WriteError(w, ErrDatabaseUnavailable, err.Error())
			return
		}

		w.Header().Set(hu.ContentType, hu.ContentTypeText)
		w.WriteHeader(http.StatusOK)
	}
}
