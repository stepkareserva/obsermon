package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stepkareserva/obsermon/internal/server/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func getTestObjects(t *testing.T) (*gomock.Controller, *mocks.MockService, *httptest.Server) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockService(ctrl)

	handlers, err := New(zap.NewNop(), mockService)
	require.NoError(t, err, "handlers initialization error")

	ts := httptest.NewServer(handlers)

	return ctrl, mockService, ts
}

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
