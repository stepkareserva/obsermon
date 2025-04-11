package handlers

import (
	"fmt"
	"net/http"
	"strings"
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

func WriteError(w http.ResponseWriter, err HandlerError, details ...string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(err.StatusCode)

	errText := fmt.Sprintln(err.Message, strings.Join(details, " "))
	w.Write([]byte(errText))
}
