package middleware

import (
	"compress/gzip"
	"fmt"
	"net/http"

	hc "github.com/stepkareserva/obsermon/internal/server/httpconst"
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

func (c *gzipWriter) Write(data []byte) (int, error) {
	if err := c.checkValidity(); err != nil {
		return 0, err
	}
	// part of http.ResponceWriter's interface contract:
	// it writes StatusOK on Write if it was not called before.
	if c.status == 0 {
		c.WriteHeader(http.StatusOK)
	}

	// compress if compressor exists
	if c.compressor != nil {
		return c.compressor.Write(data)
	} else {
		return c.ResponseWriter.Write(data)
	}
}

func (c *gzipWriter) WriteHeader(status int) {
	c.status = status

	useCompression := !c.isErrorStatus(status) &&
		c.supportContentCompress(c.Header().Get(hc.ContentType))

	if useCompression {
		c.Header().Set(hc.ContentEncoding, hc.GZipEncoding)
		c.compressor = gzip.NewWriter(c.ResponseWriter)
	} else {
		c.compressor = nil
	}

	c.ResponseWriter.WriteHeader(status)
}

func (w *gzipWriter) isErrorStatus(status int) bool {
	return status >= 400
}

func (c *gzipWriter) Close() error {
	if err := c.checkValidity(); err != nil {
		return err
	}
	if c.compressor != nil {
		return c.compressor.Close()
	}
	return nil
}

func (c *gzipWriter) supportContentCompress(contentType string) bool {
	compressableContent := []string{
		hc.ContentTypeJSON,
		hc.ContentTypeJSONU,
		hc.ContentTypeHTML,
		hc.ContentTypeHTMLU,
	}
	for _, c := range compressableContent {
		if contentType == c {
			return true
		}
	}
	return false
}

func (c *gzipWriter) checkValidity() error {
	if c == nil || c.ResponseWriter == nil {
		return fmt.Errorf("Compressed writer not exists")
	}
	return nil
}
