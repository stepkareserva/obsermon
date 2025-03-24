package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stepkareserva/obsermon/cmd/server/internal/metrics/server"
	"github.com/stepkareserva/obsermon/cmd/server/internal/metrics/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type expected struct {
	code int
}
type testCase struct {
	request  string
	expected expected
}

func TestUpdateCounter(t *testing.T) {
	updateHandler := mockUpdatesHandler(t)

	ts := httptest.NewServer(updateHandler)
	defer ts.Close()

	testCases := []testCase{
		// correct
		{request: "/counter/name/1", expected: expected{code: http.StatusOK}},
		// without metric name
		{request: "/counter/", expected: expected{code: http.StatusNotFound}},
		// incorrect metric value
		{request: "/counter/name/value", expected: expected{code: http.StatusBadRequest}},
		{request: "/counter/name/1.000", expected: expected{code: http.StatusBadRequest}},
		{request: "/counter/name/1.25", expected: expected{code: http.StatusBadRequest}},
		{request: "/counter/name/99999999999999999999999999", expected: expected{code: http.StatusBadRequest}},
	}

	for _, test := range testCases {
		t.Run(fmt.Sprintf("test '%s'", test.request), func(t *testing.T) {
			res, err := http.Post(ts.URL+test.request, "text/plain", nil)
			require.NoError(t, err)
			defer res.Body.Close()
			assert.Equal(t, test.expected.code, res.StatusCode)
		})
	}
}

func TestUpdateGauge(t *testing.T) {
	updateHandler := mockUpdatesHandler(t)

	ts := httptest.NewServer(updateHandler)
	defer ts.Close()

	testCases := []testCase{
		// correct
		{request: "/gauge/name/1.0", expected: expected{code: http.StatusOK}},
		// without metric name
		{request: "/gauge/", expected: expected{code: http.StatusNotFound}},
		// incorrect metric value
		{request: "/gauge/name/value", expected: expected{code: http.StatusBadRequest}},
		{request: "/gauge/name/1.2.3", expected: expected{code: http.StatusBadRequest}},
	}

	for _, test := range testCases {
		t.Run(fmt.Sprintf("test '%s'", test.request), func(t *testing.T) {
			res, err := http.Post(ts.URL+test.request, "text/plain", nil)
			require.NoError(t, err)
			defer res.Body.Close()
			assert.Equal(t, test.expected.code, res.StatusCode)
		})
	}
}

func mockUpdatesHandler(t *testing.T) http.Handler {
	storage := storage.NewMemStorage()
	server, err := server.NewServer(storage)
	require.NoError(t, err, "server initialization error")
	updateHandler, err := UpdateHandler(server)
	require.NoError(t, err, "update handler initialization error")
	return updateHandler
}
