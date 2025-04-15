package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// create middleware for requests and responses logging
func Logger(logger *zap.Logger) Middleware {

	return func(next http.Handler) http.Handler {
		logFn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			var responseInfo responseInfo
			next.ServeHTTP(withResponseInfo(w, &responseInfo), r)

			duration := time.Since(start)

			logger.Info("request",
				zap.String("uri", r.RequestURI),
				zap.String("method", r.Method),
				zap.Int("status", responseInfo.status),
				zap.Duration("duration", duration),
				zap.Int("size", responseInfo.size),
			)
		}

		return http.HandlerFunc(logFn)
	}
}

type responseInfo struct {
	status int
	size   int
	err    error
}

func withResponseInfo(w http.ResponseWriter, info *responseInfo) http.ResponseWriter {
	return &responseMiddleware{
		ResponseWriter: w,
		info:           info,
	}
}

type responseMiddleware struct {
	http.ResponseWriter
	info *responseInfo
}

var _ http.ResponseWriter = (*responseMiddleware)(nil)

func (m *responseMiddleware) Write(data []byte) (int, error) {
	// part of http.ResponseWriter's interface contract:
	// it writes StatusOK on Write if it was not called before.
	// due to lack of go's «inheritance» we can not intercept
	// such WriteHeader call from m.writer.Write (it calls directly)
	// m.writer.WriteHeader, not m.WriteHeader, so call it manually.
	if m.info.status == 0 {
		m.WriteHeader(http.StatusOK)
	}

	size, err := m.ResponseWriter.Write(data)
	m.info.size += size
	return size, err
}

func (m *responseMiddleware) WriteHeader(status int) {
	m.ResponseWriter.WriteHeader(status)
	m.info.status = status
}
