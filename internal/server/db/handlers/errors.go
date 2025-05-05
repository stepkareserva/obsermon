package handlers

import (
	"net/http"

	hu "github.com/stepkareserva/obsermon/internal/server/httputils"
)

var (
	ErrDatabaseUnavailable = hu.HandlerError{
		StatusCode: http.StatusInternalServerError,
		Message:    "Database unavailable",
	}
)
