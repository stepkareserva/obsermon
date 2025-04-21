package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type HandlerError struct {
	StatusCode int
	Message    string
}

var (
	ErrInvalidMetricType = HandlerError{
		StatusCode: http.StatusBadRequest,
		Message:    "Request contains invalid metric type",
	}

	ErrInvalidMetricValue = HandlerError{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid metric value",
	}

	ErrInternalServerError = HandlerError{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal server error",
	}

	ErrMissingMetricName = HandlerError{
		StatusCode: http.StatusNotFound,
		Message:    "Metric name is missing",
	}

	ErrMetricNotFound = HandlerError{
		StatusCode: http.StatusNotFound,
		Message:    "Metric not found",
	}

	ErrUnsupportedContentType = HandlerError{
		StatusCode: http.StatusUnsupportedMediaType,
		Message:    "Content-Type header is not supported",
	}

	ErrInvalidRequestJSON = HandlerError{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid request JSON content",
	}
)

func WriteError(w http.ResponseWriter, err HandlerError, log *zap.Logger, details ...string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(err.StatusCode)

	errText := fmt.Sprintln(err.Message, strings.Join(details, " "))
	if _, err := w.Write([]byte(errText)); err != nil && log != nil {
		log.Error("writing http error", zap.Error(err))
	}
}
