package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stepkareserva/obsermon/internal/server/metrics/service"
	"github.com/stepkareserva/obsermon/internal/server/metrics/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type UpdateExpected struct {
	code int
}
type UpdateTestCase struct {
	request  string
	expected UpdateExpected
}

func TestUpdateCounter(t *testing.T) {
	updateHandler := mockUpdatesHandler(t)

	ts := httptest.NewServer(updateHandler)
	defer ts.Close()

	testCases := []UpdateTestCase{
		// correct
		{request: "/counter/name/1", expected: UpdateExpected{code: http.StatusOK}},
		// without metric name
		{request: "/counter/", expected: UpdateExpected{code: http.StatusNotFound}},
		// incorrect metric value
		{request: "/counter/name/value", expected: UpdateExpected{code: http.StatusBadRequest}},
		{request: "/counter/name/1.000", expected: UpdateExpected{code: http.StatusBadRequest}},
		{request: "/counter/name/1.25", expected: UpdateExpected{code: http.StatusBadRequest}},
		{request: "/counter/name/99999999999999999999999999", expected: UpdateExpected{code: http.StatusBadRequest}},
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

	testCases := []UpdateTestCase{
		// correct
		{request: "/gauge/name/1.0", expected: UpdateExpected{code: http.StatusOK}},
		// without metric name
		{request: "/gauge/", expected: UpdateExpected{code: http.StatusNotFound}},
		// incorrect metric value
		{request: "/gauge/name/value", expected: UpdateExpected{code: http.StatusBadRequest}},
		{request: "/gauge/name/1.2.3", expected: UpdateExpected{code: http.StatusBadRequest}},
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
	service, err := service.New(storage)
	require.NoError(t, err, "service initialization error")
	updateHandler, err := UpdateHandler(service)
	require.NoError(t, err, "update handler initialization error")
	return updateHandler
}
