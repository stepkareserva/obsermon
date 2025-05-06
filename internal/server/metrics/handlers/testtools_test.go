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

func safeCloseRes(t *testing.T, res *http.Response) {
	if res == nil {
		return
	}
	err := res.Body.Close()
	require.NoError(t, err)
}
