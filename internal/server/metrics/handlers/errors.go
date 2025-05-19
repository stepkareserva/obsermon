package handlers

import (
	"net/http"
)

type HandlerError struct {
	StatusCode int
	Message    string
}

var (
	ErrInvalidMetricType = hu.HandlerError{
		StatusCode: http.StatusBadRequest,
		Message:    "Request contains invalid metric type",
	}

	ErrInvalidMetricValue = hu.HandlerError{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid metric value",
	}

	ErrMissingMetricName = hu.HandlerError{
		StatusCode: http.StatusNotFound,
		Message:    "Metric name is missing",
	}

	ErrMetricNotFound = hu.HandlerError{
		StatusCode: http.StatusNotFound,
		Message:    "Metric not found",
	}

	ErrDatabaseUnavailable = hu.HandlerError{
		StatusCode: http.StatusInternalServerError,
		Message:    "Database unavailable",
	}
)
