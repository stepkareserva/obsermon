package middleware

import (
	"compress/gzip"
	"fmt"
	"net/http"

	"github.com/stepkareserva/obsermon/internal/server/http/constants"
)

// http.ResponseWriter impl which supports
// gzip compression for some kinds of content
type gzipWriter struct {
	http.ResponseWriter
	compressor *gzip.Writer
	status     int
}

var _ http.ResponseWriter = (*gzipWriter)(nil)

func newGZipWriter(w http.ResponseWriter) (*gzipWriter, error) {
	return &gzipWriter{
		ResponseWriter: w,
	}, nil
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	if err := g.checkValidity(); err != nil {
		return 0, err
	}
	// part of http.ResponseWriter's interface contract:
	// it writes StatusOK on Write if it was not called before.
	if g.status == 0 {
		g.WriteHeader(http.StatusOK)
	}

	// compress if compressor exists
	if g.compressor != nil {
		return g.compressor.Write(data)
	}
	return g.ResponseWriter.Write(data)
}

func (g *gzipWriter) WriteHeader(status int) {
	g.status = status

	useCompression := !g.isErrorStatus(status) &&
		g.supportContentCompress(g.Header().Get(constants.ContentType))

	if useCompression {
		g.Header().Set(constants.ContentEncoding, constants.GZipEncoding)
		g.compressor = gzip.NewWriter(g.ResponseWriter)
	} else {
		g.compressor = nil
	}

	g.ResponseWriter.WriteHeader(status)
}

func (g *gzipWriter) isErrorStatus(status int) bool {
	return status >= http.StatusBadRequest
}

func (g *gzipWriter) Close() error {
	if err := g.checkValidity(); err != nil {
		return err
	}
	if g.compressor != nil {
		return g.compressor.Close()
	}
	return nil
}

func (g *gzipWriter) supportContentCompress(contentType string) bool {
	compressableContent := []string{
		constants.ContentTypeJSON,
		constants.ContentTypeJSONU,
		constants.ContentTypeHTML,
		constants.ContentTypeHTMLU,
	}
	for _, g := range compressableContent {
		if contentType == g {
			return true
		}
	}
	return false
}

func (g *gzipWriter) checkValidity() error {
	if g == nil || g.ResponseWriter == nil {
		return fmt.Errorf("compressed writer not exists")
	}
	return nil
}
