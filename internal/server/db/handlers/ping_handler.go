package handlers

import (
	"context"
	"net/http"

	hc "github.com/stepkareserva/obsermon/internal/server/httpconst"
	"go.uber.org/zap"
)

type PingHandler struct {
	db Database
}

func NewPingHandler(db Database, log *zap.Logger) (*PingHandler, error) {
	return &PingHandler{
		db: db,
	}, nil
}

func (h *PingHandler) Handler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h.db.Ping(); err != nil {
			w.Header().Set(hc.ContentType, hc.ContentTypeText)
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Header().Set(hc.ContentType, hc.ContentTypeText)
		w.WriteHeader(http.StatusOK)
	}
}
