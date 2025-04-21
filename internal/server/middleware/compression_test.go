package middleware

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRequestCompression(t *testing.T) {
	// create handler which save incoming content type header
	// and request body
	var contentEncoding string
	var requestBody []byte

	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentEncoding = r.Header.Get("Content-Encoding")
		var err error
		requestBody, err = io.ReadAll(r.Body)
		require.NoError(t, err)
		w.WriteHeader(http.StatusOK)
	})

	compressedHandler := Compression(zap.NewNop())(mockHandler)

	ts := httptest.NewServer(compressedHandler)
	defer ts.Close()

	// post gzip request
	res := testingPostGzipJSON(t, ts.URL, `{"Hello":"World"}`)
	defer func() {
		err := res.Body.Close()
		require.NoError(t, err)
	}()

	// check request unzipped
	require.Equal(t, "gzip", contentEncoding)
	require.Equal(t, `{"Hello":"World"}`, string(requestBody))
}

func TestRequestUncompressed(t *testing.T) {
	// create handler which save incoming content type header
	// and request body
	var contentEncoding string
	var requestBody []byte

	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentEncoding = r.Header.Get("Content-Encoding")
		var err error
		requestBody, err = io.ReadAll(r.Body)
		require.NoError(t, err)
		w.WriteHeader(http.StatusOK)
	})

	compressedHandler := Compression(zap.NewNop())(mockHandler)

	ts := httptest.NewServer(compressedHandler)
	defer ts.Close()

	// post gzip request
	res, err := http.Post(ts.URL, "text/plain", strings.NewReader("Hello, World"))
	require.NoError(t, err)
	defer func() {
		err := res.Body.Close()
		require.NoError(t, err)
	}()
	require.Equal(t, http.StatusOK, res.StatusCode)

	// check request unzipped
	require.Equal(t, "", contentEncoding)
	require.Equal(t, "Hello, World", string(requestBody))
}

func TestResponseCompression(t *testing.T) {
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(`{"Hello":"World"}`))
		require.NoError(t, err)
	})

	compressedHandler := Compression(zap.NewNop())(mockHandler)

	ts := httptest.NewServer(compressedHandler)
	defer ts.Close()

	// post gzip request
	res := testingPostGzipJSON(t, ts.URL, `{}`)
	defer func() {
		err := res.Body.Close()
		require.NoError(t, err)
	}()

	require.Equal(t, "gzip", res.Header.Get("Content-Encoding"))
	body := testingUngzipBody(t, res)
	require.Equal(t, `{"Hello":"World"}`, string(body))
}

func TestResponseUncompressed(t *testing.T) {
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("Hello, World"))
		require.NoError(t, err)
	})

	compressedHandler := Compression(zap.NewNop())(mockHandler)

	ts := httptest.NewServer(compressedHandler)
	defer ts.Close()

	// post gzip request
	res := testingPostGzipJSON(t, ts.URL, `{}`)
	defer func() {
		err := res.Body.Close()
		require.NoError(t, err)
	}()

	require.Equal(t, "", res.Header.Get("Content-Encoding"))
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, "Hello, World", string(body))
}

func testingPostGzipJSON(t *testing.T, url string, data string) *http.Response {
	// compress request body
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	_, err := gzipWriter.Write([]byte(data))
	require.NoError(t, err)
	err = gzipWriter.Close()
	require.NoError(t, err)

	// create request
	req, err := http.NewRequest(http.MethodPost, url, &buf)
	require.NoError(t, err)

	// set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", buf.Len()))

	// send request
	client := &http.Client{}
	res, err := client.Do(req)
	require.NoError(t, err)

	return res
}

func testingUngzipBody(t *testing.T, res *http.Response) string {
	require.Equal(t, "gzip", res.Header.Get("Content-Encoding"))
	z, err := gzip.NewReader(res.Body)
	require.NoError(t, err)
	body, err := io.ReadAll(z)
	require.NoError(t, err)
	return string(body)
}
