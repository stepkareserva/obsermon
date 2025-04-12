package middleware

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
)

// io.ReadCloser impl which supports gzip-compressed reading
type gzipReader struct {
	wrappedReader io.ReadCloser
	uncompressor  *gzip.Reader
}

var _ io.ReadCloser = (*gzipReader)(nil)

func newGZipReader(r io.ReadCloser) (*gzipReader, error) {
	uncompressor, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &gzipReader{
		wrappedReader: r,
		uncompressor:  uncompressor,
	}, nil
}

func (c *gzipReader) Read(p []byte) (n int, err error) {
	if err := c.checkValidity(); err != nil {
		return 0, err
	}
	// read from compressed reader which read from wrapped reader
	return c.uncompressor.Read(p)
}

func (c *gzipReader) Close() error {
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

func (c *gzipReader) checkValidity() error {
	if c == nil || c.uncompressor == nil || c.wrappedReader == nil {
		return fmt.Errorf("Compressed reader not exists")
	}
	return nil
}
