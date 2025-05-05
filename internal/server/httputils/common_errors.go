package httputils

import (
	"net/http"
)

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
)
