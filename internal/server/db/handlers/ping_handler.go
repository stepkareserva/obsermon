package handlers

import (
	"context"
	"net/http"

	hu "github.com/stepkareserva/obsermon/internal/server/httputils"
	"github.com/stepkareserva/obsermon/internal/server/interfaces/database"
	"go.uber.org/zap"
)

type PingHandler struct {
	db database.Database
	hu.ErrorsWriter
}

func NewPingHandler(db database.Database, log *zap.Logger) (*PingHandler, error) {
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
