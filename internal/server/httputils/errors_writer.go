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
	if writingErr := writeError(w, err); writingErr != nil {
		e.log.Error("error writing", zap.Error(writingErr))
	}

	// log error details to log
	e.log.Error("internal server error",
		zap.String("message", err.Message),
		zap.String("details", strings.Join(details, " ")),
	)
}

func writeError(w http.ResponseWriter, err HandlerError) error {
	w.Header().Set(ContentType, ContentTypeTextU)
	w.WriteHeader(err.StatusCode)

	if _, err := w.Write([]byte(err.Message)); err != nil {
		return fmt.Errorf("writing error: %w", err)
	}
	return nil
}
