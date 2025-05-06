package handlers

import (
	"context"
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/stepkareserva/obsermon/internal/server/interfaces/database"
	"go.uber.org/zap"
)

type Handlers struct {
	pingHandler *PingHandler
}

func New(db database.Database, log *zap.Logger) (*Handlers, error) {
	pingHandler, err := NewPingHandler(db, log)
	if err != nil {
		return nil, fmt.Errorf("ping handler creation: %w", err)
	}

	return &Handlers{
		pingHandler: pingHandler,
	}, nil

}

func (h *Handlers) RegisterRoutes(ctx context.Context, r chi.Router) error {
	if h == nil {
		return fmt.Errorf("handlers not exist")
	}

	r.Get("/ping", h.pingHandler.Handler(ctx))

	return nil
}
