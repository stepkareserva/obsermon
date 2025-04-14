package handlers

import (
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
