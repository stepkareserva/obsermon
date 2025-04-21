package middleware

import (
	"log"
	"net/http"

	"go.uber.org/zap/buffer"
)

// create middleware for buffered-writing responses
func Buffering() Middleware {
	return func(next http.Handler) http.Handler {
		buffering := func(w http.ResponseWriter, r *http.Request) {
			bw := withBuffering(w)
			next.ServeHTTP(bw, r)
			bw.FlushToClient()
		}

		return http.HandlerFunc(buffering)
	}
}

func withBuffering(w http.ResponseWriter) *bufferingWriter {
	return &bufferingWriter{
		ResponseWriter: w,
	}
}

type bufferingWriter struct {
	http.ResponseWriter
	buf    buffer.Buffer
	status int
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
		log.Printf("response sending error: %v", err)
	}
	w.buf.Reset()
}
