package handlers

import (
	"net/http"

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
	if writingErr := WriteError(w, err, details...); writingErr != nil {
		e.log.Error("error writing", zap.Error(writingErr))
	}
}
