package middleware

import (
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
)

// create middleware for buffered-writing responses
func Buffering(log *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		buffering := func(w http.ResponseWriter, r *http.Request) {
			bw := withBuffering(w, log)
			next.ServeHTTP(bw, r)
			bw.FlushToClient()
		}

		return http.HandlerFunc(buffering)
	}
}

func withBuffering(w http.ResponseWriter, log *zap.Logger) *bufferingWriter {
	if log == nil {
		log = zap.NewNop()
	}
	return &bufferingWriter{
		ResponseWriter: w,
		log:            log,
	}
}

type bufferingWriter struct {
	http.ResponseWriter
	buf    buffer.Buffer
	status int
	log    *zap.Logger
}

var _ http.ResponseWriter = (*bufferingWriter)(nil)

func (w *bufferingWriter) Write(data []byte) (int, error) {
	// part of http.ResponseWriter's interface contract:
	// it writes StatusOK on Write if it was not called before.
	if w.status == 0 {
		w.status = http.StatusOK
	}

	return w.buf.Write(data)
}

func (w *bufferingWriter) WriteHeader(status int) {
	w.status = status
	// clean buffer if error occurs to sending
	// only upcoming error content to client
	// without underwritten response content
	if w.isErrorStatus(status) {
		w.buf.Reset()
	}
}

func (w *bufferingWriter) isErrorStatus(status int) bool {
	return status >= http.StatusBadRequest
}

func (w *bufferingWriter) FlushToClient() {
	w.ResponseWriter.WriteHeader(w.status)
	if _, err := w.ResponseWriter.Write(w.buf.Bytes()); err != nil {
		// write error to log? pass logger here throughtout context?
		w.log.Error("response sending", zap.Error(err))
	}
	w.buf.Reset()
}
