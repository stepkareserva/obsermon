package server

import (
	"context"
	"errors"
	"net/http"
)

type Server struct {
	srv *http.Server
}

func New(addr string, handler http.Handler) *Server {
	return &Server{
		srv: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (s *Server) Start() <-chan error {
	errCh := make(chan error, 1)

	go func() {
		err := s.srv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	return errCh
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
