package middleware

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	hc "github.com/stepkareserva/obsermon/internal/server/httpconst"
)

func Compression() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		compression := func(w http.ResponseWriter, r *http.Request) {
			// handle zipped request - replace request body to unzipped
			if isCompressedRequest(r) {
				cr, err := newCompressedReader(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				r.Body = cr
				// don't move this code block out of this function!
				defer cr.Close()
			}

			// client supports compression - replace responce writer
			if compressedResponseSupported(r) {
				cw, err := newCompressedWriter(w)
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

// io.ReadCloser impl which supports gzip-compressed reading
type compressedReader struct {
	wrappedReader io.ReadCloser
	uncompressor  *gzip.Reader
}

var _ io.ReadCloser = (*compressedReader)(nil)

func newCompressedReader(r io.ReadCloser) (*compressedReader, error) {
	uncompressor, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressedReader{
		wrappedReader: r,
		uncompressor:  uncompressor,
	}, nil
}

func (c *compressedReader) Read(p []byte) (n int, err error) {
	if err := c.checkValidity(); err != nil {
		return 0, err
	}
	// read from compressed reader which read from wrapped reader
	return c.uncompressor.Read(p)
}

func (c *compressedReader) Close() error {
	if err := c.checkValidity(); err != nil {
		return err
	}
	var errs []error
	if err := c.uncompressor.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := c.wrappedReader.Close(); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func (c *compressedReader) checkValidity() error {
	if c == nil || c.uncompressor == nil || c.wrappedReader == nil {
		return fmt.Errorf("Compressed reader not exists")
	}
	return nil
}

// http.ResponseWriter impl which supports gzip compression for
// some kinds of content
type compressedWriter struct {
	http.ResponseWriter
	compressor *gzip.Writer
	status     int
}

var _ http.ResponseWriter = (*compressedWriter)(nil)

func newCompressedWriter(w http.ResponseWriter) (*compressedWriter, error) {
	return &compressedWriter{
		ResponseWriter: w,
		compressor:     gzip.NewWriter(w),
	}, nil
}

func (c *compressedWriter) Write(data []byte) (int, error) {
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

func (c *compressedWriter) WriteHeader(status int) {
	c.status = status
	if c.supportContentCompress(c.Header().Values(hc.ContentType)) {
		c.Header().Set(hc.ContentEncoding, hc.GZipEncoding)
	}
	c.ResponseWriter.WriteHeader(status)
}

func (c *compressedWriter) Close() error {
	if err := c.checkValidity(); err != nil {
		return err
	}
	return c.compressor.Close()
}

func (c *compressedWriter) supportContentCompress(contentType []string) bool {
	compressableContent := []string{
		hc.ContentTypeJSON,
		hc.ContentTypeHTML,
	}
	for _, c := range contentType {
		if lookupHeaderComponent(compressableContent, c) {
			return true
		}
	}
	return false
}

func (c *compressedWriter) checkValidity() error {
	if c == nil || c.compressor == nil || c.ResponseWriter == nil {
		return fmt.Errorf("Compressed writer not exists")
	}
	return nil
}
