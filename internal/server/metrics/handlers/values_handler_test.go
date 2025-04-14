package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ValuesExpected struct {
	code        int
	contentType string
}
type ValuesTestCase struct {
	updateRequest string
	expected      ValuesExpected
}

func TestValuesHandler(t *testing.T) {

	testCases := []ValuesTestCase{
		// correct counter
		{
			updateRequest: "/update/counter/name/2",
			expected: ValuesExpected{
				code:        http.StatusOK,
				contentType: "text/html; charset=utf-8",
			},
		},
		// correct gauge
		{
			updateRequest: "/update/gauge/name/2.5",
			expected: ValuesExpected{
				code:        http.StatusOK,
				contentType: "text/html; charset=utf-8",
			},
		},
	}

	for _, test := range testCases {
		valueHandler := mockValuesHandler(t)
		ts := httptest.NewServer(valueHandler)
		defer ts.Close()

		t.Run(fmt.Sprintf("test '%s'", test.updateRequest), func(t *testing.T) {
			// update counter
			res, err := http.Post(ts.URL+test.updateRequest, "text/plain", nil)
			require.NoError(t, err)
			defer res.Body.Close()
			require.Equal(t, http.StatusOK, res.StatusCode)

			// get values
			res, err = http.Get(ts.URL + "/")
			require.NoError(t, err)
			defer res.Body.Close()

			// check status and contentType if status is ok
			assert.Equal(t, test.expected.code, res.StatusCode)
			contentType := res.Header.Get("Content-Type")
			_, err = io.ReadAll(res.Body)
			require.NoError(t, err)
			if test.expected.code == http.StatusOK {
				assert.Equal(t, test.expected.contentType, contentType)
			}
		})
	}
}

func mockValuesHandler(t *testing.T) http.Handler {
	//storage := storage.NewMemStorage()
	//service, err := service.New(storage)
	var err error

	//storage := storage.NewMemStorage()

	//storage := storage.NewMemStorage()
	//service, err := service.New(storage)
	require.NoError(t, err, "service initialization error")
	updateHandler, err := UpdateHandler(nil)
	require.NoError(t, err, "update handler initialization error")
	valuesHandler, err := ValuesHandler(nil)
	require.NoError(t, err, "value handler initialization error")

	r := chi.NewRouter()
	r.Mount("/update", updateHandler)
	r.Mount("/", valuesHandler)

	return r
}
