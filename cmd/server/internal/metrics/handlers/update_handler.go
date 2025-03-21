package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/stepkareserva/obsermon/cmd/server/internal/metrics/server"
	"github.com/stepkareserva/obsermon/internal/models"
)

type UpdateHandler struct {
	server *server.Server
}

func NewUpdateHandler(s *server.Server) (*UpdateHandler, error) {
	if s == nil {
		return nil, fmt.Errorf("metrics server is nil")
	}

	handler := UpdateHandler{
		server: s,
	}

	return &handler, nil
}

func (h *UpdateHandler) UpdateGaugeHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.checkValidity(); err != nil {
		WriteError(w, ErrInternalServerError)
		return
	}

	var gauge models.UpdateGaugeRequest
	if err := gauge.FromURLPath(strings.TrimPrefix(r.URL.Path, "/")); err != nil {
		WriteError(w, translateURLParsingError(err))
		return
	}

	if err := h.server.UpdateGauge(gauge.Name, gauge.Value); err != nil {
		WriteError(w, ErrInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (h *UpdateHandler) UpdateCounterHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.checkValidity(); err != nil {
		WriteError(w, ErrInternalServerError)
		return
	}
	var counter models.UpdateCounterRequest
	if err := counter.FromURLPath(strings.TrimPrefix(r.URL.Path, "/")); err != nil {
		WriteError(w, translateURLParsingError(err))
		return
	}

	if err := h.server.UpdateCounter(counter.Name, counter.Value); err != nil {
		WriteError(w, ErrInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (h *UpdateHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.checkValidity(); err != nil {
		WriteError(w, ErrInternalServerError)
		return
	}
	WriteError(w, ErrInvalidMetricType)
}

func (h *UpdateHandler) checkValidity() error {
	if h == nil || h.server == nil {
		return fmt.Errorf("server not exists")
	}
	return nil
}

func translateURLParsingError(err error) Error {
	switch err.(type) {
	case *models.MissingMetricNameError:
		return ErrMissingMetricName
	case *models.InvalidMetricValueError:
		return ErrInvalidMetricValue
	default:
		return ErrInternalServerError
	}
}
