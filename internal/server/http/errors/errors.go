package errors

import (
	"net/http"
)

type HandlerError struct {
	StatusCode int
	Message    string
}

var (
	ErrInternalServerError = HandlerError{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal server error",
	}

	ErrUnsupportedContentType = HandlerError{
		StatusCode: http.StatusUnsupportedMediaType,
		Message:    "Content-Type header is not supported",
	}

	ErrInvalidRequestJSON = HandlerError{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid request JSON content",
	}

	ErrInvalidMetricType = HandlerError{
		StatusCode: http.StatusBadRequest,
		Message:    "Request contains invalid metric type",
	}

	ErrInvalidMetricValue = HandlerError{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid metric value",
	}

	ErrMissingMetricName = HandlerError{
		StatusCode: http.StatusNotFound,
		Message:    "Metric name is missing",
	}

	ErrMetricNotFound = HandlerError{
		StatusCode: http.StatusNotFound,
		Message:    "Metric not found",
	}

	ErrDatabaseUnavailable = HandlerError{
		StatusCode: http.StatusInternalServerError,
		Message:    "Database unavailable",
	}
)
