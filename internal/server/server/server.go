package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/stepkareserva/obsermon/internal/server/config"
	"github.com/stepkareserva/obsermon/internal/server/handlers"

	"github.com/stepkareserva/obsermon/internal/server/routing"
	"go.uber.org/zap"
)

type Server struct {
	service  handlers.Service
	database handlers.Database
	http     *http.Server
	log      *zap.Logger
}

func New(cfg *config.Config, s handlers.Service, db handlers.Database, log *zap.Logger) (*Server, error) {
	if s == nil {
		return nil, fmt.Errorf("service not exists")
	}
	// if db == nil {
	//		not an error, skipping
	// 		return nil, fmt.Errorf("database not exists")
	// }
	if cfg == nil {
		return nil, fmt.Errorf("config not exists")
	}
	if log == nil {
		log = zap.NewNop()
	}

	routing, err := routing.New(s, db, log)
	if err != nil {
		return nil, fmt.Errorf("handlers creator initialization: %w", err)
	}

	handler, err := routing.Handler(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("handlers initialization: %w", err)
	}

	return &Server{
		service: s,
		http: &http.Server{
			Addr:    cfg.Endpoint,
			Handler: handler,
		},
		log: log,
	}, nil
}

func (s *Server) Start() error {
	if s == nil {
		return fmt.Errorf("server not exists")
	}
	go func() {
		if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Error("HTTP server", zap.Error(err))
		}
	}()

	return nil
}

func (s *Server) Stop() error {
	// cancel server
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.http.Shutdown(context); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	return nil
}
