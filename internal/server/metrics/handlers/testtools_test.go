package handlers

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func testingGetURL(t *testing.T, url string) *http.Response {
	res, err := http.Get(url)
	require.NoError(t, err)
	return res
}

func testingPostURL(t *testing.T, url string) *http.Response {
	res, err := http.Post(url, "text/plain", nil)
	require.NoError(t, err)
	return res
}

func testingPostJSON(t *testing.T, url string, data string) *http.Response {
	res, err := http.Post(url, "application/json", strings.NewReader(data))
	require.NoError(t, err)
	return res
}

func testingGetGzipURL(t *testing.T, url string) *http.Response {
	// create request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	// set headers
	req.Header.Set("Accept-Encoding", "gzip")

	// send request
	client := &http.Client{}
	res, err := client.Do(req)
	require.NoError(t, err)

	return res
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
