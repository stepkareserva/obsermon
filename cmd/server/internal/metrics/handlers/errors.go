package handlers

import (
	"fmt"
	"net/http"
	"strings"
)

type Error struct {
	StatusCode int
	Message    string
}

var (
	ErrInvalidMetricType = Error{
		StatusCode: http.StatusBadRequest,
		Message:    "Request contains invalid metric type",
	}

	ErrInvalidMetricValue = Error{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid metric value",
	}

	ErrInternalServerError = Error{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal server error",
	}

	ErrMissingMetricName = Error{
		StatusCode: http.StatusNotFound,
		Message:    "Metric name is missing",
	}

	ErrMetricNotFound = Error{
		StatusCode: http.StatusNotFound,
		Message:    "Metric not found",
	}
)

func WriteError(w http.ResponseWriter, err Error, details ...string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(err.StatusCode)

	errText := fmt.Sprintln(err.Message, strings.Join(details, " "))
	w.Write([]byte(errText))
}
