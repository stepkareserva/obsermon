package httputils

import (
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type ErrorsWriter struct {
	log *zap.Logger
}

func NewErrorsWriter(log *zap.Logger) ErrorsWriter {
	if log == nil {
		log = zap.NewNop()
	}
	return ErrorsWriter{log: log}
}

func (e *ErrorsWriter) WriteError(w http.ResponseWriter, err HandlerError, details ...string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(err.StatusCode)

	errText := fmt.Sprintln(err.Message, strings.Join(details, " "))
	if _, err := w.Write([]byte(errText)); err != nil {
		e.log.Error("error writing", zap.Error(err))
		return
	}
}

func (e *ErrorsWriter) WriteInternalServerError(w http.ResponseWriter, details ...string) {
	e.WriteError(w, ErrInternalServerError, details...)
}

func (e *ErrorsWriter) WriteUnsupportedContentType(w http.ResponseWriter, details ...string) {
	e.WriteError(w, ErrUnsupportedContentType, details...)
}

func (e *ErrorsWriter) WriteInvalidRequestJSON(w http.ResponseWriter, details ...string) {
	e.WriteError(w, ErrInvalidRequestJSON, details...)
}
