package middleware

import (
	"net/http"
	"strings"

	hc "github.com/stepkareserva/obsermon/internal/server/httpconst"
)

func Compression() Middleware {
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
				defer cr.Close()
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
				defer cw.Close()
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(compression)
	}
}

func isCompressedRequest(r *http.Request) bool {
	return lookupHeaderComponent(
		r.Header.Values(hc.ContentEncoding),
		hc.GZipEncoding)
}

func compressedResponseSupported(r *http.Request) bool {
	return lookupHeaderComponent(
		r.Header.Values(hc.AcceptEncoding),
		hc.GZipEncoding)
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
