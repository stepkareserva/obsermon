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

func (g *gzipReader) Read(p []byte) (n int, err error) {
	if err := g.checkValidity(); err != nil {
		return 0, err
	}
	// read from compressed reader which read from wrapped reader
	return g.uncompressor.Read(p)
}

func (g *gzipReader) Close() error {
	if err := g.checkValidity(); err != nil {
		return err
	}
	var errs []error
	if err := g.uncompressor.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := g.wrappedReader.Close(); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func (g *gzipReader) checkValidity() error {
	if g == nil || g.uncompressor == nil || g.wrappedReader == nil {
		return fmt.Errorf("compressed reader not exists")
	}
	return nil
}
