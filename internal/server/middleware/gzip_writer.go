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
		compressor:     gzip.NewWriter(w),
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

	return c.compressor.Write(data)
}

func (c *gzipWriter) WriteHeader(status int) {
	c.status = status
	if c.supportContentCompress(c.Header().Get(hc.ContentType)) {
		c.Header().Set(hc.ContentEncoding, hc.GZipEncoding)
	}
	c.ResponseWriter.WriteHeader(status)
}

func (c *gzipWriter) Close() error {
	if err := c.checkValidity(); err != nil {
		return err
	}
	return c.compressor.Close()
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
	if c == nil || c.compressor == nil || c.ResponseWriter == nil {
		return fmt.Errorf("Compressed writer not exists")
	}
	return nil
}
