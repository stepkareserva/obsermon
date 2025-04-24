package middleware

import (
	"net/http"
	"strings"

	"go.uber.org/zap"

	hu "github.com/stepkareserva/obsermon/internal/server/httputils"
)

func Compression(log *zap.Logger) Middleware {
	if log == nil {
		log = zap.NewNop()
	}

	return func(next http.Handler) http.Handler {
		compression := func(w http.ResponseWriter, r *http.Request) {
			// handle zipped request - replace request body to unzipped
			if isCompressedRequest(r) {
				cr, err := newGZipReader(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				r.Body = cr
				// don't move this code block out of this function!
				defer func() {
					if err := cr.Close(); err != nil {
						log.Error("compressed reader closing", zap.Error(err))
					}
				}()
			}

			// client supports compression - replace response writer
			if compressedResponseSupported(r) {
				cw, err := newGZipWriter(w)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w = cw
				// don't move this code block out of this function!
				defer func() {
					if err := cw.Close(); err != nil {
						log.Error("compressed writer closing", zap.Error(err))
					}
				}()
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(compression)
	}
}

func isCompressedRequest(r *http.Request) bool {
	return lookupHeaderComponent(
		r.Header.Values(hu.ContentEncoding),
		hu.GZipEncoding)
}

func compressedResponseSupported(r *http.Request) bool {
	return lookupHeaderComponent(
		r.Header.Values(hu.AcceptEncoding),
		hu.GZipEncoding)
}

func lookupHeaderComponent(header []string, target string) bool {
	for _, values := range header {
		for _, value := range strings.Split(values, ",") {
			if strings.TrimSpace(value) == target {
				return true
			}
		}
	}
	return false
}
