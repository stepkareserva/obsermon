package handlers

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	hu "github.com/stepkareserva/obsermon/internal/server/httputils"
)

type DBHandler struct {
	database Database
	hu.ErrorsWriter
}

func NewDBHandler(db Database, log *zap.Logger) (*DBHandler, error) {
	return &DBHandler{
		database:     db,
		ErrorsWriter: hu.NewErrorsWriter(log),
	}, nil
}

func (h *DBHandler) PingHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h.database.Ping(); err != nil {
			h.ErrorsWriter.WriteInternalServerError(w)
		}

		w.Header().Set(hu.ContentType, hu.ContentTypeText)
		w.WriteHeader(http.StatusOK)
	}
}
