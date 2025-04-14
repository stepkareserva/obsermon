package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stepkareserva/obsermon/internal/models"
	"github.com/stepkareserva/obsermon/internal/server/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestValuesHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	valuesHandler, err := ValuesHandler(mockService)
	require.NoError(t, err, "value handler initialization error")

	ts := httptest.NewServer(valuesHandler)
	defer ts.Close()

	t.Run("test /", func(t *testing.T) {

		mockService.
			EXPECT().
			ListGauges().
			Return(models.GaugesList{}, nil)

		mockService.
			EXPECT().
			ListCounters().
			Return(models.CountersList{}, nil)

		// get values
		res, err := http.Get(ts.URL + "/")
		require.NoError(t, err)
		defer res.Body.Close()

		// check status and contentType if status is ok
		assert.Equal(t, http.StatusOK, res.StatusCode)
		contentType := res.Header.Get("Content-Type")
		_, err = io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, "text/html; charset=utf-8", contentType)
	})
}
