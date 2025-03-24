package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stepkareserva/obsermon/cmd/server/internal/metrics/server"
	"github.com/stepkareserva/obsermon/cmd/server/internal/metrics/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ValueExpected struct {
	code  int
	value string
}
type ValueTestCase struct {
	updateRequest string
	valueRequest  string
	expected      ValueExpected
}

func TestValueHandler(t *testing.T) {

	testCases := []ValueTestCase{
		// correct counter
		{
			updateRequest: "/update/counter/name/2",
			valueRequest:  "/value/counter/name",
			expected: ValueExpected{
				code:  http.StatusOK,
				value: "2",
			},
		},
		// correct gauge
		{
			updateRequest: "/update/gauge/name/2.5",
			valueRequest:  "/value/gauge/name",
			expected: ValueExpected{
				code:  http.StatusOK,
				value: "2.5",
			},
		},
		// incorrect counter
		{
			updateRequest: "/update/counter/name/2",
			valueRequest:  "/value/counter/othername",
			expected: ValueExpected{
				code:  http.StatusNotFound,
				value: "",
			},
		},
	}

	for _, test := range testCases {
		valueHandler := mockValueHandler(t)
		ts := httptest.NewServer(valueHandler)
		defer ts.Close()

		t.Run(fmt.Sprintf("test '%s'", test.valueRequest), func(t *testing.T) {
			// update counter
			res, err := http.Post(ts.URL+test.updateRequest, "text/plain", nil)
			require.NoError(t, err)
			defer res.Body.Close()
			require.Equal(t, http.StatusOK, res.StatusCode)

			// get counter value
			res, err = http.Get(ts.URL + test.valueRequest)
			require.NoError(t, err)
			defer res.Body.Close()

			// check status and responce if status is ok
			assert.Equal(t, test.expected.code, res.StatusCode)
			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			if test.expected.code == http.StatusOK {
				assert.Equal(t, test.expected.value, string(body))
			}
		})
	}
}

func mockValueHandler(t *testing.T) http.Handler {
	storage := storage.NewMemStorage()
	server, err := server.NewServer(storage)
	require.NoError(t, err, "server initialization error")
	updateHandler, err := UpdateHandler(server)
	require.NoError(t, err, "update handler initialization error")
	valueHandler, err := ValueHandler(server)
	require.NoError(t, err, "value handler initialization error")

	r := chi.NewRouter()
	r.Mount("/update", updateHandler)
	r.Mount("/value", valueHandler)

	return r
}
