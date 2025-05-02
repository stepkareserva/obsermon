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
	if writingErr := writeError(w, err, details...); writingErr != nil {
		e.log.Error("error writing", zap.Error(writingErr))
	}
}

func writeError(w http.ResponseWriter, err HandlerError, details ...string) error {
	w.Header().Set(ContentType, ContentTypeTextU)
	w.WriteHeader(err.StatusCode)

	errText := fmt.Sprintln(err.Message, strings.Join(details, " "))
	if _, err := w.Write([]byte(errText)); err != nil {
		return fmt.Errorf("writing error: %w", err)
	}
	return nil
}
